-- Seed default categories (idempotent: ON CONFLICT DO NOTHING)
INSERT INTO categories (id, name, description, slug, parent_id, "order", icon, created_at)
VALUES
  ('general', 'General Discussion', 'General topics and community chat', 'general', NULL, 0, '', NOW()),
  ('installation', 'Installation & Setup', 'Installation issues and setup help', 'installation', NULL, 1, '', NOW()),
  ('features', 'Feature Requests', 'Suggestions and feature requests', 'features', NULL, 2, '', NOW()),
  ('development', 'Development', 'Development and contributing', 'development', NULL, 3, '', NOW())
ON CONFLICT (id) DO NOTHING;
