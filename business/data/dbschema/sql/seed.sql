INSERT INTO users (user_id, name, email, roles, password_hash, date_created, date_updated) VALUES
	('5cf37266-3473-4006-984f-9325122678b7', 'AFP Admin', 'admin@afp.com', '{ADMIN,USER}', '$2a$10$1ggfMVZV6Js0ybvJufLRUOWHS5f6KneuP0XwwHpJ8L8ipdry9f2/a', '2022-03-24 00:00:00', '2022-03-24 00:00:00'),
	('45b5fbd3-755f-4379-8f07-a58d4a30fa2f', 'Getty User', 'user@getty.com', '{USER}', '$2a$10$9/XASPKBbJKVfCAZKDH.UuhsuALDr5vVm6VrYA9VFR8rccK86C1hW', '2022-03-24 00:00:00', '2022-03-24 00:00:00')
	ON CONFLICT DO NOTHING;

INSERT INTO images (image_id, image_url, user_id, date_uploaded) VALUES
	('a2b0639f-2cc6-44b8-b97b-15d69dbb511e', 'images/afp.jpg', '5cf37266-3473-4006-984f-9325122678b7', '2022-01-01 00:00:01.000001+00'),
	('72f8b983-3eb4-48db-9ed0-e45cc6bd716b', 'images/getty.jpg', '5cf37266-3473-4006-984f-9325122678b7', '2022-01-01 00:00:02.000001+00')
	ON CONFLICT DO NOTHING;