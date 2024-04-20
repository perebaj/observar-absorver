CREATE TABLE essays (
    title TEXT NOT NULL,
    body TEXT NOT NULL,
    body_embedding vector(1536) NOT NULL,
    model_name VARCHAR(200) NOT NULL
);
