package services

import (
	"bytes"
	"context"
	"errors"
	"html/template"
	"io/fs"
	"os"
	"sync"

	"github.com/harmonify/movie-reservation-system/notification-service/internal/core/shared"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type (
	EmailTemplateService interface {
		// Render renders HTML email template
		Render(ctx context.Context, path shared.EmailTemplatePath, data interface{}) (string, error)
	}

	EmailTemplateServiceParam struct {
		fx.In
		fx.Lifecycle

		EmailTemplatePaths []shared.EmailTemplatePath `group:"email-template-paths"`
		Logger             logger.Logger
		Tracer             tracer.Tracer
	}

	EmailTemplateServiceResult struct {
		fx.Out

		EmailTemplateService EmailTemplateService
	}

	emailTemplateServiceImpl struct {
		cache  sync.Map
		logger logger.Logger
		tracer tracer.Tracer
	}
)

func NewEmailTemplateService(p EmailTemplateServiceParam) EmailTemplateServiceResult {
	s := &emailTemplateServiceImpl{
		cache:  sync.Map{},
		logger: p.Logger,
		tracer: p.Tracer,
	}

	p.Lifecycle.Append(fx.StartHook(func(ctx context.Context) error {
		return s.preloadTemplates(ctx, p.EmailTemplatePaths)
	}))

	return EmailTemplateServiceResult{
		EmailTemplateService: s,
	}
}

func (s *emailTemplateServiceImpl) Render(ctx context.Context, path shared.EmailTemplatePath, data interface{}) (string, error) {
	ctx, span := s.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	content, err := s.getTemplate(ctx, path)
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to get email template", zap.Error(err), zap.Any("template_path", path))
		return "", err
	}
	var body bytes.Buffer
	if err := content.Execute(&body, data); err != nil {
		s.logger.WithCtx(ctx).Error("Failed to render email template", zap.Error(err), zap.Any("template_path", path), zap.Any("template_data", data))
		return "", err
	}
	return body.String(), nil
}

// getTemplate fetches and parses the template from the cache or local file
func (s *emailTemplateServiceImpl) getTemplate(ctx context.Context, path shared.EmailTemplatePath) (*template.Template, error) {
	ctx, span := s.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	// Load from cache, if possible
	if content, found := s.cache.Load(path); found {
		return content.(*template.Template), nil
	}
	// Load from local file
	content, err := s.loadTemplate(ctx, path)
	if err != nil {
		return nil, err
	}
	// Store to the cache
	s.cache.Store(path, content)
	return content, nil
}

// loadTemplate loads a template from a file and parses it
func (s *emailTemplateServiceImpl) loadTemplate(ctx context.Context, path shared.EmailTemplatePath) (*template.Template, error) {
	ctx, span := s.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	if _, err := os.Stat(path.String()); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			s.logger.WithCtx(ctx).Error("Template file does not exist", zap.String("path", path.String()))
		} else {
			s.logger.WithCtx(ctx).Error("Unknown error while reading template file", zap.Error(err), zap.String("path", path.String()))
		}
		return nil, err
	}
	tmplFile, err := os.ReadFile(path.String())
	if err != nil {
		return nil, err
	}
	tmpl, err := template.New(path.String()).Parse(string(tmplFile))
	if err != nil {
		return nil, err
	}
	return tmpl, nil
}

// preloadTemplates cache all templates to the memory
func (s *emailTemplateServiceImpl) preloadTemplates(ctx context.Context, paths []shared.EmailTemplatePath) error {
	for _, path := range paths {
		content, err := s.loadTemplate(ctx, path)
		if err != nil {
			return err
		}
		s.cache.Store(path, content)
	}
	return nil
}
