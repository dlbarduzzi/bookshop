ALTER TABLE books
ADD CONSTRAINT authors_length_check
CHECK (array_length(authors, 1) BETWEEN 1 AND 3);

ALTER TABLE books
ADD CONSTRAINT books_published_date_check
CHECK (published_date >= '1000-01-01' AND published_date < now());

ALTER TABLE books
ADD CONSTRAINT books_page_count_check
CHECK (page_count >= 1);

ALTER TABLE books
ADD CONSTRAINT categories_length_check
CHECK (array_length(categories, 1) BETWEEN 1 AND 5);
