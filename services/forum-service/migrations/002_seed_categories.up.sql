INSERT INTO categories (id, name, description, slug, "order") VALUES
('cat-general', 'General', 'General discussions', 'general', 1),
('cat-news', 'News', 'Project news and announcements', 'news', 2),
('cat-guides', 'Guides & Tutorials', 'Community guides', 'guides', 3)
ON CONFLICT DO NOTHING;
