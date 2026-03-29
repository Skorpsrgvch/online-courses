import React from 'react';
import './HomePage.css';

const HomePage = () => {
  return (
    <div className="home-page">
      

      {/* Главный экран (Hero) */}
      <section className="hero">
        <div className="container hero-content">
          <h2>Эксперт по женскому здоровью</h2>
          <p>Помогаю женщинам обрести гармонию, здоровье и уверенность в себе через комплексный подход.</p>
          <div className="hero-buttons">
            <button className="btn-primary">Записаться на консультацию</button>
            <button className="btn-secondary">Узнать больше</button>
          </div>
        </div>
      </section>

      {/* О специалисте */}
      <section id="about" className="section about">
        <div className="container">
          <h2>О специалисте</h2>
          <div className="achievements">
            <div className="achievement-item">🏆 10+ лет опыта</div>
            <div className="achievement-item">🎓 5 международных сертификатов</div>
            <div className="achievement-item">❤️ 1000+ довольных пациенток</div>
          </div>
        </div>
      </section>

      {/* Услуги */}
      <section id="services" className="section services">
        <div className="container">
          <h2>Услуги</h2>
          <div className="cards-grid">
            <div className="card">
              <h3>Консультация</h3>
              <p>Индивидуальный разбор вашей ситуации</p>
            </div>
            <div className="card">
              <h3>Ведение</h3>
              <p>Полное сопровождение до результата</p>
            </div>
            <div className="card">
              <h3>Онлайн-курсы</h3>
              <p>Обучающие программы для самостоятельной работы</p>
            </div>
          </div>
        </div>
      </section>

      {/* Курсы (заглушка) */}
      <section id="courses" className="section courses">
        <div className="container">
          <h2>Популярные курсы</h2>
          <p>Здесь будет список курсов из базы данных</p>
        </div>
      </section>

      {/* Подвал */}
      <footer className="footer">
        <div className="container">
          <p>© 2024 Женское Здоровье. Все права защищены.</p>
        </div>
      </footer>
    </div>
  );
};

export default HomePage;