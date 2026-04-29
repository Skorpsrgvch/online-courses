import React from 'react';
import { sanitizeArticleHtml } from '../../lib/dompurify';

interface ArticleViewerProps {
  content: string;
}

export const ArticleViewer: React.FC<ArticleViewerProps> = ({ content }) => {
  const safeHtml = sanitizeArticleHtml(content);

  return (
    <div
      className="prose prose-gray max-w-none article-content"
      dangerouslySetInnerHTML={{ __html: safeHtml }}
    />
  );
};
