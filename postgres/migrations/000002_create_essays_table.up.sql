CREATE TABLE essays (
    ID TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    url TEXT NOT NULL,
    content TEXT NOT NULL,
    date TEXT NOT NULL,
    embedding vector(1536) NOT NULL,
    model_name VARCHAR(200) NOT NULL,
    dimension INT NOT NULL
);
