import { useMemo, useState } from 'react';
import { Link } from 'react-router-dom';
import {
  ChevronLeftIcon,
  ChevronRightIcon,
} from '@heroicons/react/24/outline';

const QUALIFICATIONS_COUNT = 31;

const QualificationsPage = () => {
  const images = useMemo(
    () =>
      Array.from({ length: QUALIFICATIONS_COUNT }, (_, index) => ({
        id: index + 1,
        src: `/images/Qualifications/${index + 1}.jpg`,
        alt: `Документ о квалификации ${index + 1}`,
      })),
    []
  );

  const [activeIndex, setActiveIndex] = useState(0);
  const [touchStart, setTouchStart] = useState<number | null>(null);

  const activeImage = images[activeIndex];

  const goToPrevious = () => {
    setActiveIndex((current) =>
      current === 0 ? images.length - 1 : current - 1
    );
  };

  const goToNext = () => {
    setActiveIndex((current) =>
      current === images.length - 1 ? 0 : current + 1
    );
  };

  const handleTouchEnd = (clientX: number) => {
    if (touchStart === null) return;

    const diff = touchStart - clientX;
    const minSwipeDistance = 50;

    if (Math.abs(diff) >= minSwipeDistance) {
      if (diff > 0) {
        goToNext();
      } else {
        goToPrevious();
      }
    }

    setTouchStart(null);
  };

  return (
    <main className="min-h-screen bg-gray-50 py-8 sm:py-10 lg:py-14">
        
      <section className="mx-auto w-full max-w-7xl px-4 sm:px-6 lg:px-8">
        <Link to="/" className="inline-flex items-start text-sm md:text-base lg:text-base font-medium text-gray-500! hover:text-rose-500! mb-4 transition-colors"
                  style={{ textDecoration: 'none' }}>
                    ← Назад на главную
                  </Link>
        <div className="mb-6 text-center sm:mb-8">
          <h1 className="text-3xl font-semibold text-gray-900 sm:text-4xl lg:text-5xl">
            Образование и сертификаты
          </h1>
          <p className="mx-auto mt-3 max-w-2xl text-sm leading-6 text-gray-600 sm:text-base">
            Документы, подтверждающие профессиональное обучение, практику и повышение квалификации.
          </p>
        </div>

        <div className="relative mx-auto max-w-5xl">
          <button
            type="button"
            onClick={goToPrevious}
            aria-label="Предыдущий документ"
            className="absolute left-5 top-1/2 z-20 flex h-8 w-8 -translate-y-1/2 items-center justify-center rounded-full bg-white/85 text-gray-700 shadow-md shadow-gray-900/10 ring-1 ring-rose-100 transition hover:bg-white hover:text-rose-600 focus:outline-none focus:ring-2 focus:ring-rose-300 sm:left-6 sm:h-9 sm:w-9"
          >
            <ChevronLeftIcon className="h-4 w-4 sm:h-5 sm:w-5" aria-hidden="true" />
          </button>

          <div
            className="overflow-hidden rounded-[28px] bg-white p-3 shadow-xl shadow-rose-950/10 ring-1 ring-rose-100 sm:p-4 lg:p-5"
            onTouchStart={(event) => setTouchStart(event.touches[0].clientX)}
            onTouchEnd={(event) => handleTouchEnd(event.changedTouches[0].clientX)}
          >
            <div className="flex min-h-[420px] items-center justify-center rounded-[22px] bg-gradient-to-br from-stone-50 via-white to-rose-50 p-3 sm:min-h-[620px] sm:p-5 lg:min-h-[720px]">
              <img
                key={activeImage.id}
                src={activeImage.src}
                alt={activeImage.alt}
                className="max-h-[70vh] w-full max-w-full select-none object-contain drop-shadow-sm"
                loading={activeIndex === 0 ? 'eager' : 'lazy'}
                draggable={false}
              />
            </div>
          </div>

          <button
            type="button"
            onClick={goToNext}
            aria-label="Следующий документ"
            className="absolute right-5 top-1/2 z-20 flex h-8 w-8 -translate-y-1/2 items-center justify-center rounded-full bg-white/85 text-gray-700 shadow-md shadow-gray-900/10 ring-1 ring-rose-100 transition hover:bg-white hover:text-rose-600 focus:outline-none focus:ring-2 focus:ring-rose-300 sm:right-6 sm:h-9 sm:w-9"
          >
            <ChevronRightIcon className="h-4 w-4 sm:h-5 sm:w-5" aria-hidden="true" />
          </button>
        </div>

        <div className="mx-auto mt-5 flex max-w-5xl flex-col items-center gap-4 sm:mt-6">
          <div className="rounded-full bg-white px-4 py-2 text-sm font-medium text-gray-700 shadow-sm ring-1 ring-rose-100">
            {activeIndex + 1} / {images.length}
          </div>

          <div className="w-full overflow-x-auto pb-2">
            <div className="mx-auto flex w-max gap-2 px-1 sm:gap-3">
              {images.map((image, index) => {
                const isActive = index === activeIndex;

                return (
                  <button
                    key={image.id}
                    type="button"
                    onClick={() => setActiveIndex(index)}
                    aria-label={`Открыть документ ${image.id}`}
                    className={`h-16 w-12 shrink-0 overflow-hidden rounded-lg bg-white p-0.5 ring-2 transition sm:h-20 sm:w-14 ${
                      isActive
                        ? 'ring-rose-500 shadow-md shadow-rose-900/10'
                        : 'ring-transparent opacity-70 hover:opacity-100 hover:ring-rose-200'
                    }`}
                  >
                    <img
                      src={image.src}
                      alt=""
                      className="h-full w-full rounded-md object-cover"
                      loading="lazy"
                    />
                  </button>
                );
              })}
            </div>
          </div>
        </div>
      </section>
    </main>
  );
};

export default QualificationsPage;
