import React from 'react';
import { Link } from 'react-router-dom';

interface PolicyPlaceholderPageProps {
  title: string;
  text: React.ReactNode;
}

const PolicyPlaceholderPage: React.FC<PolicyPlaceholderPageProps> = ({ title, text }) => {
  return (
    <div className="min-h-screen bg-gray-50 text-gray-900">
      <main className="mx-auto flex min-h-[70vh] max-w-6xl flex-col justify-center px-4 py-12 sm:px-8 lg:px-10">
        <Link
          to="/"
          className="mb-10 inline-flex w-fit items-center text-sm font-medium text-gray-600! transition-colors hover:text-rose-600!"
          style={{ textDecoration: 'none' }}
        >
          ← На главную
        </Link>

        <section className="rounded-2xl border border-gray-200 bg-white p-8 shadow-md sm:p-12">
          <p className="mb-4! text-sm font-semibold uppercase tracking-wider text-rose-500">Документ</p>
          <h1 className="mb-8! font-serif text-4xl font-bold! text-gray-900 sm:text-5xl">{title}</h1>
          <div className="prose prose-lg max-w-none text-gray-800 leading-relaxed">
            {text}
          </div>
        </section>
      </main>
    </div>
  );
};

export default PolicyPlaceholderPage;
