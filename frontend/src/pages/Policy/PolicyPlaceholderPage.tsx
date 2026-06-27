import React from 'react';
import { Link } from 'react-router-dom';

interface PolicyPlaceholderPageProps {
  title: string;
  text: string;
}

const PolicyPlaceholderPage: React.FC<PolicyPlaceholderPageProps> = ({ title, text }) => {
  return (
    <div className="min-h-screen bg-gray-50 text-gray-900">
      <main className="mx-auto flex min-h-[70vh] max-w-4xl flex-col justify-center px-4 py-12 sm:px-6 lg:px-8">
        <Link
          to="/"
          className="mb-8 inline-flex w-fit items-center text-sm font-medium text-gray-500! transition-colors hover:text-rose-600!"
          style={{ textDecoration: 'none' }}
        >
          ← На главную
        </Link>

        <section className="rounded-2xl border border-gray-100 bg-white p-6 shadow-sm sm:p-10">
          <p className="mb-3 text-sm font-medium uppercase tracking-wide text-rose-500">Документ</p>
          <h1 className="mb-6 font-serif text-3xl font-bold text-gray-900 sm:text-4xl">{title}</h1>
          <p className="text-lg leading-relaxed text-gray-600">{text}</p>
        </section>
      </main>
    </div>
  );
};

export default PolicyPlaceholderPage;
