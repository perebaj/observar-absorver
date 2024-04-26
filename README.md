# observar-absorver

Use the public resources of Eduardo Marinho Philosopher to create a GPT model that can generate text in the same style as him.

# Get Starting

All available commands can be found using `make help`.

# Services

Services are components of code that have specific functions. Between them, we have the following services:

- Scraper: This service is responsible for scraping the data from the website of Eduardo Marinho.
- GPT: This service is responsible for generating embeddings upon the text data.
- Snippets: This service is responsible for generating snippets of text from the embeddings. (Split in smaller parts)

# Environment Variables

**Obligatory:**

- OPENAI_API_KEY: The API key for the OpenAI service.
