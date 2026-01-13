-- Migration: Add FULLTEXT search index for topics
-- Date: 2026-01-13
-- Description: Enables natural language search on topic titles with relevance ranking

ALTER TABLE topics ADD FULLTEXT INDEX idx_title_fulltext (title);
