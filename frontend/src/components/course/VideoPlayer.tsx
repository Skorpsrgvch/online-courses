import React, { useEffect, useMemo } from 'react';

interface VideoPlayerProps {
  url: string;
  onProgress?: (percent: number) => void;
}

const normalizeRuTubeEmbedUrl = (url: string) => {
  if (!url) return '';

  try {
    const parsedUrl = new URL(url);

    if (parsedUrl.hostname.includes('rutube.ru')) {
      parsedUrl.pathname = parsedUrl.pathname.replace('/video/embed/', '/play/embed/');
    }

    return parsedUrl.toString();
  } catch {
    return url;
  }
};

export const VideoPlayer: React.FC<VideoPlayerProps> = ({ url, onProgress }) => {
  const embedUrl = useMemo(() => normalizeRuTubeEmbedUrl(url), [url]);

  useEffect(() => {
    if (!onProgress) return;

    const checkProgress = setInterval(() => {
      onProgress(10);
    }, 5000);

    return () => clearInterval(checkProgress);
  }, [onProgress]);

  return (
    <div className="relative h-[60vw] max-h-[600px] min-h-[240px] w-full overflow-hidden rounded-xl bg-black shadow-lg sm:aspect-video sm:h-auto sm:min-h-0">
      <iframe
        src={embedUrl}
        title="Video Lesson"
        className="absolute inset-0 h-full w-full border-0"
        allow="autoplay; fullscreen; picture-in-picture; encrypted-media; clipboard-write; screen-wake-lock"
        allowFullScreen
        loading="eager"
        {...{ webkitAllowFullScreen: true, mozallowfullscreen: 'true' }}
      />
    </div>
  );
};
