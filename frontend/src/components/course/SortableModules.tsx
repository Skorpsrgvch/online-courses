// src/components/course/SortableCourseStructure.tsx
import React from 'react';
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

// --- Компонент перетаскиваемого урока ---
interface SortableLessonProps {
    lesson: LessonInput;
    moduleIndex: number;
    lessonIndex: number;
    updateLesson: (mIdx: number, lIdx: number, field: keyof LessonInput, value: any) => void;
    removeLesson: (mIdx: number, lIdx: number) => void;
}

const SortableLesson: React.FC<SortableLessonProps> = ({
    lesson,
    moduleIndex,
    lessonIndex,
    updateLesson,
    removeLesson,
}) => {
    const { attributes, listeners, setNodeRef, transform, transition, isDragging } = useSortable({
        id: `lesson-${lesson.id}-${lessonIndex}`, // Уникальный ID для DnD
    });

    const style = {
        transform: CSS.Transform.toString(transform),
        transition,
        opacity: isDragging ? 0.5 : 1,
        zIndex: isDragging ? 10 : 1,
    };

    return (
        <div ref={setNodeRef} style={style} className="p-4 border-b last:border-0 grid grid-cols-12 gap-3 items-start bg-white">
            {/* Ручка для перетаскивания */}
            <div className="col-span-1 flex items-center justify-center pt-6 cursor-grab active:cursor-grabbing text-gray-400 hover:text-gray-600" {...attributes} {...listeners}>
                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 8h16M4 16h16" /></svg>
            </div>

            <div className="col-span-12 md:col-span-2">
                <label className="block text-xs text-gray-500 mb-1">Название урока</label>
                <input
                    type="text"
                    value={lesson.title}
                    onChange={(e) => updateLesson(moduleIndex, lessonIndex, 'title', e.target.value)}
                    className="w-full text-sm px-2 py-1 border rounded focus:ring-rose-500"
                    placeholder="Название"
                />
            </div>
            <div className="col-span-12 md:col-span-2">
                <label className="block text-xs text-gray-500 mb-1">Video Embed ID</label>
                <input
                    type="text"
                    value={lesson.video_embed_id}
                    onChange={(e) => updateLesson(moduleIndex, lessonIndex, 'video_embed_id', e.target.value)}
                    className="w-full text-sm px-2 py-1 border rounded focus:ring-rose-500"
                    placeholder="YouTube ID"
                />
            </div>
            <div className="col-span-12 md:col-span-2">
                <label className="block text-xs text-gray-500 mb-1">Private Key</label>
                <input
                    type="text"
                    value={lesson.private_key || ''}
                    onChange={(e) => updateLesson(moduleIndex, lessonIndex, 'private_key', e.target.value)}
                    className="w-full text-sm px-2 py-1 border rounded focus:ring-rose-500"
                    placeholder="Secret"
                />
            </div>
            <div className="col-span-12 md:col-span-2 flex items-end">
                <button type="button" onClick={() => removeLesson(moduleIndex, lessonIndex)} className="text-red-500 hover:text-red-700 text-xs font-medium underline">
                    Удалить
                </button>
            </div>

            <div className="col-span-12 md:col-span-5">
                <label className="block text-xs text-gray-500 mb-1">Описание</label>
                <textarea
                    value={lesson.description}
                    onChange={(e) => updateLesson(moduleIndex, lessonIndex, 'description', e.target.value)}
                    rows={1}
                    className="w-full text-sm px-2 py-1 border rounded focus:ring-rose-500"
                    placeholder="Кратко..."
                />
            </div>
        </div>
    );
};

// --- Компонент перетаскиваемого модуля ---
interface SortableModuleProps {
    module: ModuleInput;
    mIndex: number;
    updateModule: (idx: number, field: keyof ModuleInput, value: any) => void;
    removeModule: (idx: number) => void;
    addLesson: (mIdx: number) => void;
    updateLesson: (mIdx: number, lIdx: number, field: keyof LessonInput, value: any) => void;
    removeLesson: (mIdx: number, lIdx: number) => void;
    // Для вложенного DnD уроков нам нужно знать порядок уроков внутри этого модуля
    // Но так как уроки тоже сортируемые, мы передадим весь массив уроков и функцию их перемещения
    onLessonsReorder: (mIdx: number, oldIndex: number, newIndex: number) => void;
}

const SortableModule: React.FC<SortableModuleProps> = ({
    module,
    mIndex,
    updateModule,
    removeModule,
    addLesson,
    updateLesson,
    removeLesson,
    onLessonsReorder,
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

    // Обработчик окончания перетаскивания для уроков внутри этого модуля
    const handleLessonsDragEnd = (event: DragEndEvent) => {
        const { active, over } = event;
        if (over && active.id !== over.id) {
            // Извлекаем индексы из ID вида "lesson-id-index"
            const activeIdStr = active.id.toString();
            const overIdStr = over.id.toString();
            
            const activeIndex = parseInt(activeIdStr.split('-').pop() || '0');
            const overIndex = parseInt(overIdStr.split('-').pop() || '0');

            onLessonsReorder(mIndex, activeIndex, overIndex);
        }
    };

    // Сенсоры для вложенного DnD
    const sensors = useSensors(
        useSensor(PointerSensor),
        useSensor(KeyboardSensor, { coordinateGetter: sortableKeyboardCoordinates })
    );

    return (
        <div ref={setNodeRef} style={style} className="border border-gray-200 rounded-xl bg-gray-50 p-4 mb-4 shadow-sm">
            {/* Заголовок модуля с ручкой */}
            <div className="flex gap-4 items-start mb-4">
                <div className="pt-1 cursor-grab active:cursor-grabbing text-gray-400 hover:text-gray-600" {...attributes} {...listeners}>
                    <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 8h16M4 16h16" /></svg>
                </div>
                <div className="flex-1">
                    <label className="block text-xs font-bold text-gray-500 uppercase mb-1">Название модуля</label>
                    <input
                        type="text"
                        value={module.title}
                        onChange={(e) => updateModule(mIndex, 'title', e.target.value)}
                        placeholder="Введите название модуля"
                        className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-rose-500 bg-white font-semibold"
                    />
                </div>
                <div className="pt-6">
                    <button type="button" onClick={() => removeModule(mIndex)} className="text-red-500 hover:text-red-700 text-sm font-medium px-3 py-1 border border-red-200 rounded hover:bg-red-50">
                        Удалить
                    </button>
                </div>
            </div>

            {/* Список уроков с собственным DnD контекстом */}
            <div className="bg-white rounded-lg border border-gray-200 overflow-hidden ml-10">
                <div className="bg-gray-100 px-4 py-2 text-xs font-bold text-gray-600 uppercase flex justify-between">
                    <span>Уроки (перетаскивайте для смены порядка)</span>
                    <span>Действия</span>
                </div>

                {module.lessons.length === 0 && (
                    <div className="p-4 text-center text-sm text-gray-500 italic">В этом модуле пока нет уроков</div>
                )}

                <DndContext
                    sensors={sensors}
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
                    <button type="button" onClick={() => addLesson(mIndex)} className="text-sm text-rose-600 font-medium hover:text-rose-800 flex items-center gap-1">
                        + Добавить урок
                    </button>
                </div>
            </div>
        </div>
    );
};

// --- Главный контейнер структуры курса ---
interface CourseStructureEditorProps {
    modules: ModuleInput[];
    setModules: React.Dispatch<React.SetStateAction<ModuleInput[]>>;
    updateModule: (idx: number, field: keyof ModuleInput, value: any) => void;
    removeModule: (idx: number) => void;
    addLesson: (mIdx: number) => void;
    updateLesson: (mIdx: number, lIdx: number, field: keyof LessonInput, value: any) => void;
    removeLesson: (mIdx: number, lIdx: number) => void;
    addModule: () => void;
}

export const CourseStructureEditor: React.FC<CourseStructureEditorProps> = ({
    modules,
    setModules,
    updateModule,
    removeModule,
    addLesson,
    updateLesson,
    removeLesson,
    addModule,
}) => {
    // Обработчик перетаскивания модулей
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

    // Обработчик перетаскивания уроков (обновление состояния в родителе)
    const handleLessonsReorder = (mIdx: number, oldIndex: number, newIndex: number) => {
        setModules((prevModules) => {
            const newModules = [...prevModules];
            const currentModule = { ...newModules[mIdx] };
            currentModule.lessons = arrayMove(currentModule.lessons, oldIndex, newIndex);
            newModules[mIdx] = currentModule;
            return newModules;
        });
    };

    const sensors = useSensors(
        useSensor(PointerSensor, { activationConstraint: { distance: 5 } }), // Небольшая задержка чтобы не мешать кликам
        useSensor(KeyboardSensor, { coordinateGetter: sortableKeyboardCoordinates })
    );

    return (
        <div className="bg-white rounded-2xl shadow-sm border border-gray-100 overflow-hidden">
            <div className="p-6 sm:p-8">
                <div className="flex justify-between items-center mb-6 border-b pb-2">
                    <h2 className="text-xl font-bold text-gray-800">Программа курса</h2>
                    <button type="button" onClick={addModule} className="px-4 py-2 bg-green-50 text-green-700 font-medium rounded-lg hover:bg-green-100 transition-colors">
                        + Добавить модуль
                    </button>
                </div>

                <DndContext
                    sensors={sensors}
                    collisionDetection={closestCenter}
                    onDragEnd={handleModulesDragEnd}
                >
                    <SortableContext items={modules.map((_, idx) => `module-${idx}`)} strategy={verticalListSortingStrategy}>
                        <div className="space-y-2">
                            {modules.map((module, mIndex) => (
                                <SortableModule
                                    key={`module-${module.id}-${mIndex}`}
                                    module={module}
                                    mIndex={mIndex}
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
    );
};