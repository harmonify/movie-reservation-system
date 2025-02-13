db["movies"].createIndex(
    {
        title: "text",
        description: "text",
        genres: "text",
        "cast.name": "text",
        "director.name": "text",
        "writer.name": "text",
        production_company: "text",
    },
    {
        weights: {
            title: 10,
            description: 3,
            genres: 5,
            "cast.name": 3,
            "director.name": 2,
            "writer.name": 2,
            production_company: 1,
        },
    }
);
