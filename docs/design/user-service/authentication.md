# Authentication

This document describes the implementation of authentication process, JWT, and algorithms used in the user service.

## Register

When a user registers on the app, the service first checks if the user already exists in the database. If the user already exists, the service returns an error message.

If the user does not exist, the service hashes the user's password using the Argon2id algorithm and stores the hashed password in the database.

> For more information on the Argon2id algorithm implementation, see the [Argon2id](#argon2id) section.

The service then generates a unique RSA key pair using 2048 bits. Both the public and private key is then encoded in PKCS#1 ASN.1 PEM format.

> For more information on PKCS#1 ASN.1 PEM format, see the [PKCS#1 ASN.1 PEM](#pkcs1-asn1-pem-format) section.

Before storing into the database, the service encode the public key in base64 format and encrypts the private key.

The service uses the AES-256 GCM with PBKDF2 key derivation algorithm to encrypt the private key.

> For more information on AES-256 GCM with PBKDF2, see the [AES-256 GCM with PBKDF2](#aes-256-gcm-with-pbkdf2) section.

The service then assign the user role and stores all the user's metadata in the database.

If the user registration is successful, the service generate a random token which is then used to construct the user's account verification link.

This token is then stored in the cache which will be used to verify the user's email.

The service sends the account verification link to the user's email and returns a success message.

### Verify email

When a user clicks on the account verification link, the service first checks if the user exists in the database.

If the user exists, the service then checks if the token in the account verification link is valid.

If the token is valid, the service updates the user metadata in the database, acknowledging the email is verified, and then returns a success message.

### Login

When a user logs in, the service first find the matching user in the database.

If the user exists, the service compares the hashed password in the database with the password provided by the user.

If the passwords match, the service signs a JWT using the RS256 algorithm using the private key that is stored securely in the database. This access token is only valid for 15 minutes.

> For more information on RS256, see the [RS256 (RSA signature with SHA-256)](#rs256-rsa-signature-with-sha-256) section.

The service also generates a refresh token using random secure 32 bytes. The refresh token can be used to generate a new JWT when the current JWT expires.

The refresh token will be hashed before storing it into the database. This ensures that the refresh token is secure even if the database is compromised.

The service then returns the access token and access token validity duration in the response body, and sets the unhashed refresh token in the response header as a HTTPS only cookie.

### Refresh token

When the user's access token expires, the user can use the refresh token to generate a new access token.

The service first hash the refresh token and check if the hashed token exists in the database.

If the hashed token exists, the service then generates a new access token and returns it in the response body.

### Logout

When the user logs out, the service first hash the refresh token and check if the hashed token exists in the database.

If the hashed token exists, the service then deletes the hashed token from the database, deletes the refresh token cookie, and returns a success message.

## Miscellaneous

### Argon2id

Argon2id is a password-hashing function that is resistant to side-channel attacks and GPU cracking attacks. It is the winner of the Password Hashing Competition (PHC) in 2015.

### Argon2id implementation

The following configuration is used for the Argon2id algorithm implementation in the user service:

-   Memory: 64 \* 1024
-   Iterations: 1
-   Parallelism: Number of CPUs
-   SaltLength: 16
-   KeyLength: 32

See <https://en.wikipedia.org/wiki/Argon2> for more.

### RSA

See <https://en.wikipedia.org/wiki/RSA_(cryptosystem)>.

#### PKCS#1 ASN.1 PEM format

A PKCS#1 ASN.1 PEM format refers to a way of:

-   storing a PKCS#1 cryptographic key (typically RSA) using the ASN.1 data structure standard,
-   encoded in the Distinguished Encoding Rules (DER) format, and then
-   wrapped within a PEM (Privacy Enhanced Mail) container.

##### ASN.1

A standard language for defining data structures, used here to specify the components of a PKCS#1 key like modulus, public exponent, and private exponent.

##### DER Encoding

A specific way of serializing ASN.1 data into a binary format, ensuring consistent representation across different systems.

##### PEM Container

PEM files are essentially base64 encoded versions of the DER encoded data.

A text-based format that includes header and footer lines like "-----BEGIN RSA PRIVATE KEY-----" and "-----END RSA PRIVATE KEY-----" to identify the key type, and then encodes the DER-encoded key data in Base64.

See <https://mbed-tls.readthedocs.io/en/latest/kb/cryptography/asn1-key-structures-in-der-and-pem/> for more.

### AES-256 GCM with PBKDF2

#### AES-256 GCM

AES-256 (Advanced Encryption Standard) in Galois/Counter Mode (GCM) provides both encryption and authentication. GCM ensures the integrity of the encrypted data and prevents tampering.

See <https://en.wikipedia.org/wiki/Galois/Counter_Mode>.

#### PBKDF2

PBKDF2 (Password-Based Key Derivation Function 2) is a key derivation function that uses a pseudorandom function to derive a strong encryption key from a password.

See <https://en.wikipedia.org/wiki/PBKDF2>.

#### AES-256 GCM with PBKDF2 implementation

PBKDF2 is used to derive a strong encryption key from the application's secret key (AppSecret) and a secure random salt.

This ensures that the derived encryption key is unique for every encryption operation, even when the same application secret is used.

In this case, the service uses this method to generate a unique key to create AES-256 cipher block wrapped in GCM.

The service then uses the cipher and a secure random nonce (number that's only used once) to encrypt user private keys.

After the encryption, the service encode the encrypted private key, the nonce, and the salt in base64 format and arrange them in the following format:

```txt
<cipher base64>.<salt base64>.<nonce base64>.<pbkdf2 iterations>
```

To decrypt the private key, the service extracts the cipher, salt, nonce, and PBKDF2 iterations from the arranged base64-encoded string and uses them to derive the encryption key.

### SHA-256

SHA256 or Secure Hash Algorithm 256, is a hashing algorithm that converts text of any length into a fixed-size string of 256 bits.

See <https://en.wikipedia.org/wiki/SHA-2>.

### RS256 (RSA signature with SHA-256)

RS256 is a cryptographic algorithm used to sign JWTs (JSON Web Tokens) using an RSA private key and verify the signature using the corresponding RSA public key.

See <https://tools.ietf.org/html/rfc7518#section-3.3>.
