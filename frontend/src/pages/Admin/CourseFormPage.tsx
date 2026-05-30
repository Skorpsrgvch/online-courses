import React, { useState, useEffect } from 'react';
import { useParams, useNavigate, Link } from 'react-router-dom';
// Импортируем типы отдельно
import type { Course, BonusItem } from '../../api/types';
import { coursesService } from '../../api/courses.service';

// DnD импорты
import {
    DndContext,
    closestCenter,
    KeyboardSensor,
    PointerSensor,
    useSensor,
    useSensors,
    type DragEndEvent,
} from '@dnd-kit/core';
import {
    arrayMove,
    SortableContext,
    sortableKeyboardCoordinates,
    verticalListSortingStrategy,
    useSortable,
} from '@dnd-kit/sortable';
import { CSS } from '@dnd-kit/utilities';

// --- Типы ---
interface LessonInput {
    id: number;
    title: string;
    description: string;
    video_embed_id: string;
    private_key: string | null;
    order: number;
}

interface ModuleInput {
    id: number;
    title: string;
    order: number;
    lessons: LessonInput[];
}

type FormMode = 'create' | 'edit';

interface CourseFormPageProps {
    mode: FormMode;
}

// --- Компоненты для сортировки ---

// 1. Сортируемый урок
const SortableLesson: React.FC<{
    lesson: LessonInput;
    moduleIndex: number;
    lessonIndex: number;
    updateLesson: (mIdx: number, lIdx: number, field: keyof LessonInput, value: any) => void;
    removeLesson: (mIdx: number, lIdx: number) => void;
}> = ({ lesson, moduleIndex, lessonIndex, updateLesson, removeLesson }) => {
    const { attributes, listeners, setNodeRef, transform, transition, isDragging } = useSortable({
        id: `lesson-${lesson.id}-${lessonIndex}`,
    });

    const style = {
        transform: CSS.Transform.toString(transform),
        transition,
        opacity: isDragging ? 0.5 : 1,
        zIndex: isDragging ? 10 : 1,
        position: 'relative' as const,
    };

    return (
        <div ref={setNodeRef} style={style} className="p-4 border-b last:border-0 grid grid-cols-12 gap-3 items-start bg-white group hover:bg-gray-50 transition-colors relative">
            {/* Ручка для перетаскивания */}
            <div className="col-span-1 flex items-center justify-center pt-6 cursor-grab active:cursor-grabbing text-gray-400 hover:text-gray-600" {...attributes} {...listeners}>
                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 8h16M4 16h16" />
                </svg>
            </div>

            <div className="col-span-12 md:col-span-2">
                <label className="block text-xs text-gray-500 mb-1">Название урока</label>
                <input
                    type="text"
                    value={lesson.title}
                    onChange={(e) => updateLesson(moduleIndex, lessonIndex, 'title', e.target.value)}
                    className="w-full text-sm px-2 py-1 border rounded focus:ring-rose-500 focus:border-rose-500"
                    placeholder="Название"
                />
            </div>
            <div className="col-span-12 md:col-span-2">
                <label className="block text-xs text-gray-500 mb-1">Video Embed ID</label>
                <input
                    type="text"
                    value={lesson.video_embed_id}
                    onChange={(e) => updateLesson(moduleIndex, lessonIndex, 'video_embed_id', e.target.value)}
                    className="w-full text-sm px-2 py-1 border rounded focus:ring-rose-500 focus:border-rose-500"
                    placeholder="YouTube ID"
                />
            </div>
            <div className="col-span-12 md:col-span-2">
                <label className="block text-xs text-gray-500 mb-1">Private Key</label>
                <input
                    type="text"
                    value={lesson.private_key || ''}
                    onChange={(e) => updateLesson(moduleIndex, lessonIndex, 'private_key', e.target.value)}
                    className="w-full text-sm px-2 py-1 border rounded focus:ring-rose-500 focus:border-rose-500"
                    placeholder="Secret"
                />
            </div>
            
            {/* Кнопка удаления - теперь видна всегда при наведении на строку урока */}
            <div className="col-span-12 md:col-span-2 flex items-end pb-1">
                <button 
                    type="button" 
                    onClick={() => removeLesson(moduleIndex, lessonIndex)} 
                    className="text-red-500 hover:text-red-700 hover:bg-red-50 p-1.5 rounded transition-colors"
                    title="Удалить урок"
                >
                    <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                    </svg>
                </button>
            </div>

            <div className="col-span-12 md:col-span-5">
                <label className="block text-xs text-gray-500 mb-1">Описание</label>
                <textarea
                    value={lesson.description}
                    onChange={(e) => updateLesson(moduleIndex, lessonIndex, 'description', e.target.value)}
                    rows={1}
                    className="w-full text-sm px-2 py-1 border rounded focus:ring-rose-500 focus:border-rose-500"
                    placeholder="Кратко..."
                />
            </div>
        </div>
    );
};

// 2. Сортируемый модуль с функцией сворачивания
const SortableModule: React.FC<{
    module: ModuleInput;
    mIndex: number;
    isExpanded: boolean;
    toggleExpand: (id: number) => void;
    updateModule: (idx: number, field: keyof ModuleInput, value: any) => void;
    removeModule: (idx: number) => void;
    addLesson: (mIdx: number) => void;
    updateLesson: (mIdx: number, lIdx: number, field: keyof LessonInput, value: any) => void;
    removeLesson: (mIdx: number, lIdx: number) => void;
    onLessonsReorder: (mIdx: number, oldIndex: number, newIndex: number) => void;
}> = ({ 
    module, 
    mIndex, 
    isExpanded, 
    toggleExpand,
    updateModule, 
    removeModule, 
    addLesson, 
    updateLesson, 
    removeLesson, 
    onLessonsReorder 
}) => {
    
    const { attributes, listeners, setNodeRef, transform, transition, isDragging } = useSortable({
        id: `module-${module.id}-${mIndex}`,
    });

    const style = {
        transform: CSS.Transform.toString(transform),
        transition,
        opacity: isDragging ? 0.5 : 1,
        zIndex: isDragging ? 20 : 1,
        position: 'relative' as const,
    };

    // Сенсоры для вложенного DnD (уроки)
    const lessonSensors = useSensors(
        useSensor(PointerSensor),
        useSensor(KeyboardSensor, { coordinateGetter: sortableKeyboardCoordinates })
    );

    const handleLessonsDragEnd = (event: DragEndEvent) => {
        const { active, over } = event;
        if (over && active.id !== over.id) {
            const activeIdStr = active.id.toString();
            const overIdStr = over.id.toString();
            const activeIndex = parseInt(activeIdStr.split('-').pop() || '0');
            const overIndex = parseInt(overIdStr.split('-').pop() || '0');
            onLessonsReorder(mIndex, activeIndex, overIndex);
        }
    };

    return (
        <div ref={setNodeRef} style={style} className={`border rounded-xl bg-gray-50 mb-4 shadow-sm overflow-hidden ${isDragging ? 'border-rose-300 ring-2 ring-rose-100' : 'border-gray-200'}`}>
            
            {/* ЗАГОЛОВОК МОДУЛЯ (Всегда виден) */}
            <div className="flex items-center p-4 bg-white border-b border-gray-100">
                {/* Ручка для перетаскивания всего модуля */}
                <div className="mr-3 cursor-grab active:cursor-grabbing text-gray-400 hover:text-gray-600 p-1" {...attributes} {...listeners}>
                    <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 8h16M4 16h16" />
                    </svg>
                </div>

                {/* Кнопка разворота/сворота */}
                <button
                    type="button"
                    onClick={() => toggleExpand(module.id)}
                    className="mr-3 p-1 rounded-full hover:bg-gray-100 text-gray-500 transition-transform duration-200"
                    style={{ transform: isExpanded ? 'rotate(90deg)' : 'rotate(0deg)' }}
                    title={isExpanded ? "Свернуть" : "Развернуть"}
                >
                    <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                    </svg>
                </button>

                {/* Поле названия */}
                <div className="flex-1 mr-4">
                    <label className="block text-xs font-bold text-gray-400 uppercase mb-0.5">Модуль {mIndex + 1}</label>
                    <input
                        type="text"
                        value={module.title}
                        onChange={(e) => updateModule(mIndex, 'title', e.target.value)}
                        placeholder="Введите название модуля"
                        className="w-full font-semibold text-gray-800 border-b border-transparent hover:border-gray-200 focus:border-rose-500 focus:ring-0 bg-transparent px-0 py-1 transition-colors"
                    />
                </div>

                {/* Статистика уроков (если свернут) */}
                {!isExpanded && (
                    <div className="mr-4 text-xs text-gray-500 bg-gray-100 px-2 py-1 rounded-full">
                        {module.lessons.length} ур.
                    </div>
                )}

                {/* Кнопка удаления */}
                <button 
                    type="button" 
                    onClick={() => removeModule(mIndex)} 
                    className="text-red-400 hover:text-red-600 hover:bg-red-50 p-2 rounded-lg transition-colors"
                    title="Удалить модуль"
                >
                    <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                    </svg>
                </button>
            </div>

            {/* ТЕЛО МОДУЛЯ (Уроки) - Показывается только если развернут */}
            {isExpanded && (
                <div className="bg-gray-50 animate-fadeIn">
                    <div className="bg-gray-100 px-4 py-2 text-xs font-bold text-gray-600 uppercase flex justify-between border-b border-gray-200">
                        <span>Уроки (перетаскивайте для смены порядка)</span>
                        <span>Действия</span>
                    </div>

                    {module.lessons.length === 0 && (
                        <div className="p-8 text-center text-sm text-gray-500 italic bg-white">
                            В этом модуле пока нет уроков. Нажмите кнопку ниже, чтобы добавить первый.
                        </div>
                    )}

                    {/* Вложенный DnD контекст для уроков */}
                    <DndContext
                        sensors={lessonSensors}
                        collisionDetection={closestCenter}
                        onDragEnd={handleLessonsDragEnd}
                    >
                        <SortableContext items={module.lessons.map((_, idx) => `lesson-${module.id}-${idx}`)} strategy={verticalListSortingStrategy}>
                            {module.lessons.map((lesson, lIndex) => (
                                <SortableLesson
                                    key={`lesson-${lesson.id}-${lIndex}`}
                                    lesson={lesson}
                                    moduleIndex={mIndex}
                                    lessonIndex={lIndex}
                                    updateLesson={updateLesson}
                                    removeLesson={removeLesson}
                                />
                            ))}
                        </SortableContext>
                    </DndContext>

                    <div className="p-3 bg-gray-50 border-t border-gray-200">
                        <button type="button" onClick={() => addLesson(mIndex)} className="text-sm text-rose-600 font-medium hover:text-rose-800 flex items-center gap-1 px-2 py-1 rounded hover:bg-rose-50 transition-colors">
                            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" /></svg>
                            Добавить урок
                        </button>
                    </div>
                </div>
            )}
        </div>
    );
};


// --- Основной компонент страницы ---

const CourseFormPage: React.FC<CourseFormPageProps> = ({ mode }) => {
    const { id } = useParams<{ id: string }>();
    const navigate = useNavigate();

    const [isLoading, setIsLoading] = useState(mode === 'edit');
    const [isSaving, setIsSaving] = useState(false);
    const [error, setError] = useState<string | null>(null);

    // Состояние основных данных курса
    const [formData, setFormData] = useState({
        title: '',
        description: '',
        is_public: false,
        price: 0,
        cover_image_url: '',
        contraindications: '',
        recommendations: '',
        target_audience: '',
        course_basis: '',
        class_basis: '',
        bonuses: [] as BonusItem[],
        is_active: true,
    });

    // Состояние модулей и уроков
    const [modules, setModules] = useState<ModuleInput[]>([]);

    // Состояние развернутых модулей (хранит ID развернутых модулей)
    const [expandedModules, setExpandedModules] = useState<Set<number>>(new Set());

    // Сенсоры для основного DnD (модули)
    const sensors = useSensors(
        useSensor(PointerSensor, { activationConstraint: { distance: 5 } }),
        useSensor(KeyboardSensor, { coordinateGetter: sortableKeyboardCoordinates })
    );

    // Загрузка данных при редактировании
    useEffect(() => {
        if (mode === 'edit' && id) {
            const loadCourse = async () => {
                try {
                    setIsLoading(true);
                    setError(null);
                    const data = await coursesService.getCourseFull(Number(id));

                    const c = data.course;
                    setFormData({
                        title: c.title,
                        description: c.description,
                        is_public: c.is_public,
                        price: c.price,
                        cover_image_url: c.cover_image_url || '',
                        contraindications: c.contraindications || '',
                        recommendations: c.recommendations || '',
                        target_audience: c.target_audience || '',
                        course_basis: c.course_basis || '',
                        class_basis: c.class_basis || '',
                        bonuses: c.bonuses || [],
                        is_active: c.is_active,
                    });

                    if (data.modules && Array.isArray(data.modules)) {
                        const formattedModules: ModuleInput[] = data.modules.map((m: any) => ({
                            id: m.id,
                            title: m.title,
                            order: m.order,
                            lessons: (m.lessons || []).map((l: any) => ({
                                id: l.id,
                                title: l.title,
                                description: l.description || '',
                                video_embed_id: l.video_embed_id || '',
                                private_key: l.private_key || null,
                                order: l.order,
                            })),
                        }));
                        formattedModules.sort((a, b) => a.order - b.order);
                        setModules(formattedModules);
                        
                        // Разворачиваем первый модуль по умолчанию для удобства
                        if (formattedModules.length > 0) {
                            setExpandedModules(new Set([formattedModules[0].id]));
                        }
                    }
                } catch (err: any) {
                    setError(err.message || 'Не удалось загрузить данные курса');
                } finally {
                    setIsLoading(false);
                }
            };
            loadCourse();
        } else {
            setIsLoading(false);
            setModules([
                { id: 0, title: 'Новый модуль', order: 1, lessons: [] },
            ]);
            // Для нового курса сразу разворачиваем единственный модуль
            setExpandedModules(new Set([0])); 
        }
    }, [mode, id]);

    const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) => {
        const { name, value, type } = e.target;
        if (type === 'checkbox') {
            const checked = (e.target as HTMLInputElement).checked;
            setFormData(prev => ({ ...prev, [name]: checked }));
            return;
        }
        if (name === 'price') {
            setFormData(prev => ({ ...prev, [name]: Number(value) }));
            return;
        }
        setFormData(prev => ({ ...prev, [name]: value }));
    };


    const addModule = () => {
        const newOrder = modules.length > 0 ? Math.max(...modules.map(m => m.order)) + 1 : 1;
        
        setModules([...modules, { 
            id: 0, 
            title: '', 
            order: newOrder, 
            lessons: [] 
        }]);
        
        
        const tempId = 0; 

        setModules(prev => {
            const newMods = [...prev];
            newMods[newMods.length - 1] = { ...newMods[newMods.length - 1], id: tempId };
            return newMods;
        });

        setExpandedModules(prev => new Set(prev).add(tempId));
    };

    const removeModule = (index: number) => {
        if (window.confirm('Вы уверены? Это удалит модуль и все его уроки.')) {
            const moduleToRemove = modules[index];
            const newModules = [...modules];
            newModules.splice(index, 1);
            setModules(newModules);
            
            setExpandedModules(prev => {
                const next = new Set(prev);
                next.delete(moduleToRemove.id);
                return next;
            });
        }
    };

    const updateModule = (index: number, field: keyof ModuleInput, value: any) => {
        const newModules = [...modules];
        newModules[index] = { ...newModules[index], [field]: value };
        setModules(newModules);
    };

    // Переключатель видимости модуля
    const toggleModuleExpand = (moduleId: number) => {
        setExpandedModules(prev => {
            const next = new Set(prev);
            if (next.has(moduleId)) {
                next.delete(moduleId);
            } else {
                next.add(moduleId);
            }
            return next;
        });
    };

    // --- Логика управления уроками ---
    const addLesson = (moduleIndex: number) => {
        const module = modules[moduleIndex];
        const newOrder = module.lessons.length > 0 ? Math.max(...module.lessons.map(l => l.order)) + 1 : 1;
        
        const newLesson: LessonInput = { 
            id: 0, // Бэкенд поймет, что это новый урок
            title: '', 
            description: '', 
            video_embed_id: '', 
            private_key: null, 
            order: newOrder 
        };
        
        const newModules = [...modules];
        newModules[moduleIndex] = { ...module, lessons: [...module.lessons, newLesson] };
        setModules(newModules);
    };

    const removeLesson = (moduleIndex: number, lessonIndex: number) => {
        const newModules = [...modules];
        newModules[moduleIndex].lessons.splice(lessonIndex, 1);
        setModules(newModules);
    };

    const updateLesson = (moduleIndex: number, lessonIndex: number, field: keyof LessonInput, value: any) => {
        const newModules = [...modules];
        const lesson = { ...newModules[moduleIndex].lessons[lessonIndex], [field]: value };
        if (field === 'private_key' && value === '') {
            lesson.private_key = null;
        }
        newModules[moduleIndex].lessons[lessonIndex] = lesson;
        setModules(newModules);
    };

    // --- Обработчики DnD ---
    
    // Перемещение модулей
    const handleModulesDragEnd = (event: DragEndEvent) => {
        const { active, over } = event;
        if (over && active.id !== over.id) {
            const activeIdStr = active.id.toString();
            const overIdStr = over.id.toString();
            const activeIndex = parseInt(activeIdStr.split('-').pop() || '0');
            const overIndex = parseInt(overIdStr.split('-').pop() || '0');
            setModules((items) => arrayMove(items, activeIndex, overIndex));
        }
    };

    // Перемещение уроков внутри модуля
    const handleLessonsReorder = (mIdx: number, oldIndex: number, newIndex: number) => {
        setModules((prevModules) => {
            const newModules = [...prevModules];
            const currentModule = { ...newModules[mIdx] };
            currentModule.lessons = arrayMove(currentModule.lessons, oldIndex, newIndex);
            newModules[mIdx] = currentModule;
            return newModules;
        });
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!formData.title.trim()) {
            setError('Название курса обязательно.');
            return;
        }

        setIsSaving(true);
        setError(null);

        try {
            const payload = {
                ...formData,
                modules: modules.map((m, idx) => ({
                    ...m,
                    order: idx + 1, 
                    lessons: m.lessons.map((l, lIdx) => ({
                        ...l,
                        order: lIdx + 1
                    }))
                }))
            };

            if (mode === 'create') {
                await coursesService.createCourseWithModules(payload); 
                navigate('/admin');
            } else {
                if (!id) throw new Error('ID курса не найден');
                await coursesService.updateFullCourse(Number(id), payload);
                navigate('/admin');
            }
        } catch (err: any) {
            console.error(err);
            setError(err.message || 'Неизвестная ошибка');
        } finally {
            setIsSaving(false);
        }
    };

    if (isLoading) {
        return (
            <div className="min-h-screen bg-gray-50 flex items-center justify-center">
                <div className="text-center">
                    <div className="w-12 h-12 border-4 border-rose-300 border-t-rose-500 rounded-full animate-spin mx-auto mb-4"></div>
                    <p className="text-gray-500">Загрузка данных курса...</p>
                </div>
            </div>
        );
    }

    // Конфигурация полей с русскими названиями
    const contentFields = [
        { key: 'target_audience', label: 'Курс для вас, если...', placeholder: 'Проблема 1|||Проблема 2' },
        { key: 'course_basis', label: 'Курс включает в себя', placeholder: 'Тема 1|||Тема 2' },
        { key: 'class_basis', label: 'Основа занятий', placeholder: 'Метод 1|||Метод 2' },
        { key: 'recommendations', label: 'Рекомендации', placeholder: '✅ Совет 1|||✅ Совет 2' },
        { key: 'contraindications', label: 'Противопоказания', placeholder: '❌ Противопоказание 1' },
    ];

    return (
        <div className="min-h-screen bg-gray-50 py-8">
            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
                <div className="mb-8">
                    <Link to="/admin" className="inline-flex items-center text-sm text-gray-500! hover:text-rose-500! mb-4 transition-colors" style={{ textDecoration: 'none' }}>
                        ← Назад в админ-панель
                    </Link>
                    <h1 className="text-3xl font-serif font-bold text-gray-900">
                        {mode === 'create' ? 'Создать новый курс' : 'Редактировать курс'}
                    </h1>
                </div>

                {error && (
                    <div className="mb-6 p-4 bg-red-50 text-red-700 rounded-lg border border-red-100">
                        {error}
                    </div>
                )}

                <form onSubmit={handleSubmit} className="space-y-8">
                    
                    {/* 1. Основная информация */}
                    <div className="bg-white rounded-2xl shadow-sm border border-gray-100 overflow-hidden">
                        <div className="p-6 sm:p-8 space-y-6">
                            <h2 className="text-xl font-bold text-gray-800 border-b pb-2">Основная информация</h2>
                            
                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-1">Название курса *</label>
                                <input type="text" name="title" value={formData.title} onChange={handleInputChange} required className="w-full px-4 py-2 border rounded-lg focus:ring-2 focus:ring-rose-500" />
                            </div>

                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-1">Описание *</label>
                                <textarea name="description" value={formData.description} onChange={handleInputChange} required rows={4} className="w-full px-4 py-2 border rounded-lg focus:ring-2 focus:ring-rose-500" />
                            </div>

                            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                                <div>
                                    <label className="block text-sm font-medium text-gray-700 mb-1">Цена (₽)</label>
                                    <input type="number" name="price" value={formData.price} onChange={handleInputChange} min="0" className="w-full px-4 py-2 border rounded-lg" />
                                </div>
                                <div className="flex items-center gap-2 mt-6">
                                    <input type="checkbox" name="is_public" checked={formData.is_public} onChange={handleInputChange} className="w-5 h-5 text-rose-600 rounded" />
                                    <span className="text-sm font-medium text-gray-700">Публичный (бесплатный)</span>
                                </div>
                                <div className="flex items-center gap-2 mt-6">
                                    <input type="checkbox" name="is_active" checked={formData.is_active} onChange={handleInputChange} className="w-5 h-5 text-rose-600 rounded" />
                                    <span className="text-sm font-medium text-gray-700">Активен</span>
                                </div>
                            </div>

                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-1">URL обложки</label>
                                <input type="text" name="cover_image_url" value={formData.cover_image_url} onChange={handleInputChange} className="w-full px-4 py-2 border rounded-lg" />
                                {formData.cover_image_url && (
                                    <img src={formData.cover_image_url} alt="Preview" className="mt-2 h-48 w-auto rounded-lg object-cover border" onError={(e) => (e.currentTarget.style.display = 'none')} />
                                )}
                            </div>
                        </div>
                    </div>

                    {/* 2. Контентные блоки (Переведены на русский) */}
                    <div className="bg-white rounded-2xl shadow-sm border border-gray-100 overflow-hidden">
                        <div className="p-6 sm:p-8 space-y-4">
                            <h2 className="text-xl font-bold text-gray-800 border-b pb-2">Детали курса</h2>
                            <p className="text-xs text-gray-500">Разделяйте пункты символом <code>|||</code></p>
                            
                            {contentFields.map((field) => (
                                <div key={field.key}>
                                    <label className="block text-sm font-medium text-gray-700 mb-1">{field.label}</label>
                                    <textarea 
                                        name={field.key} 
                                        value={formData[field.key as keyof typeof formData] as string} 
                                        onChange={handleInputChange} 
                                        rows={2} 
                                        placeholder={field.placeholder}
                                        className="w-full px-4 py-2 border rounded-lg font-mono text-sm focus:ring-rose-500 focus:border-rose-500" 
                                    />
                                </div>
                            ))}
                        </div>
                    </div>

                    {/* 3. Редактор Модулей и Уроков (DnD + Accordion) */}
                    <div className="bg-white rounded-2xl shadow-sm border border-gray-100 overflow-hidden">
                        <div className="p-6 sm:p-8">
                            <div className="flex justify-between items-center mb-6 border-b pb-2">
                                <div>
                                    <h2 className="text-xl font-bold text-gray-800">Программа курса</h2>
                                    <p className="text-xs text-gray-500 mt-1">Перетаскивайте модули за иконку ≡. Разворачивайте модули для редактирования уроков.</p>
                                </div>
                                <button type="button" onClick={addModule} className="px-4 py-2 bg-green-50 text-green-700 font-medium rounded-lg hover:bg-green-100 transition-colors flex items-center gap-2">
                                    <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" /></svg>
                                    Добавить модуль
                                </button>
                            </div>

                            <DndContext
                                sensors={sensors}
                                collisionDetection={closestCenter}
                                onDragEnd={handleModulesDragEnd}
                            >
                                <SortableContext items={modules.map((_, idx) => `module-${idx}`)} strategy={verticalListSortingStrategy}>
                                    <div className="space-y-3">
                                        {modules.map((module, mIndex) => (
                                            <SortableModule
                                                key={`module-${module.id}-${mIndex}`}
                                                module={module}
                                                mIndex={mIndex}
                                                isExpanded={expandedModules.has(module.id)}
                                                toggleExpand={toggleModuleExpand}
                                                updateModule={updateModule}
                                                removeModule={removeModule}
                                                addLesson={addLesson}
                                                updateLesson={updateLesson}
                                                removeLesson={removeLesson}
                                                onLessonsReorder={handleLessonsReorder}
                                            />
                                        ))}
                                    </div>
                                </SortableContext>
                            </DndContext>
                        </div>
                    </div>

                    {/* 4. Бонусы (JSON) */}
                    <div className="bg-white rounded-2xl shadow-sm border border-gray-100 overflow-hidden">
                        <div className="p-6 sm:p-8">
                            <h2 className="text-xl font-bold text-gray-800 mb-4 border-b pb-2">Бонусы (JSON)</h2>
                            <textarea
                                value={JSON.stringify(formData.bonuses, null, 2)}
                                onChange={(e) => {
                                    try {
                                        const parsed = JSON.parse(e.target.value);
                                        setFormData(prev => ({ ...prev, bonuses: parsed }));
                                    } catch { /* ignore */ }
                                }}
                                rows={6}
                                className="w-full px-4 py-2 border rounded-lg font-mono text-xs bg-gray-50"
                            />
                        </div>
                    </div>

                    {/* Футер */}
                    <div className="flex justify-end gap-4 pb-12">
                        <button type="button" onClick={() => navigate('/admin')} className="px-6 py-3 text-gray-700 font-medium hover:bg-gray-200 rounded-lg transition-colors">
                            Отмена
                        </button>
                        <button
                            type="submit"
                            disabled={isSaving}
                            className={`px-8 py-3 bg-rose-500 text-white font-bold rounded-lg shadow-md hover:bg-rose-600 transition-all ${isSaving ? 'opacity-70 cursor-not-allowed' : ''}`}
                        >
                            {isSaving ? 'Сохранение...' : (mode === 'create' ? 'Создать курс' : 'Сохранить все изменения')}
                        </button>
                    </div>
                </form>
            </div>
        </div>
    );
};

export default CourseFormPage;
