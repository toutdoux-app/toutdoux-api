CREATE TABLE IF NOT EXISTS "schema_migration" (
"version" TEXT NOT NULL
);
CREATE UNIQUE INDEX "schema_migration_version_idx" ON "schema_migration" (version);
CREATE TABLE IF NOT EXISTS "users" (
"id" TEXT PRIMARY KEY,
"email" TEXT NOT NULL,
"password_hash" TEXT NOT NULL,
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL
);
CREATE TABLE IF NOT EXISTS "todo_lists" (
"id" TEXT PRIMARY KEY,
"name" TEXT NOT NULL,
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL,
"user_id" char(36) NOT NULL,
FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE cascade
);
CREATE TABLE IF NOT EXISTS "todo_entries" (
"id" TEXT PRIMARY KEY,
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL,
"title" TEXT NOT NULL,
"priority" INTEGER NOT NULL,
"done" bool NOT NULL DEFAULT 'false',
"todo_list_id" char(36) NOT NULL,
FOREIGN KEY (todo_list_id) REFERENCES todo_lists (id) ON DELETE cascade
);
CREATE TABLE IF NOT EXISTS "todo_list_labels" (
"id" TEXT PRIMARY KEY,
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL,
"name" TEXT NOT NULL,
"todo_list_id" char(36) NOT NULL,
FOREIGN KEY (todo_list_id) REFERENCES todo_lists (id) ON DELETE cascade
);
CREATE TABLE IF NOT EXISTS "todo_entry_labels" (
"id" TEXT PRIMARY KEY,
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL,
"todo_entry_id" char(36) NOT NULL,
"todo_list_label_id" char(36) NOT NULL,
FOREIGN KEY (todo_entry_id) REFERENCES todo_entries (id) ON DELETE cascade,
FOREIGN KEY (todo_list_label_id) REFERENCES todo_list_labels (id) ON DELETE cascade
);
CREATE TABLE IF NOT EXISTS "todo_entry_relations" (
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL,
"relation_type" TEXT NOT NULL,
"todo_entry_id" char(36) NOT NULL,
"related_to_todo_entry_id" char(36) NOT NULL,
FOREIGN KEY (todo_entry_id) REFERENCES todo_entries (id) ON DELETE cascade,
FOREIGN KEY (related_to_todo_entry_id) REFERENCES todo_entries (id) ON DELETE cascade,
PRIMARY KEY("todo_entry_id", "related_to_todo_entry_id")
);
