import apiClient from './axiosInstance';
import type { StudentStat, StudentDetailsResponse } from './types';

export const adminService = {
    /**
     * Получить топ студентов (лидеров по прогрессу)
     */
    getTopStudents: async (limit: number = 20): Promise<StudentStat[]> => {
        const response = await apiClient.get<{ students: StudentStat[] }>('/admin/students', { params: { limit } });
        return response.data.students;
    },

    /**
     * Получить детальную информацию о студенте и его курсах
     */
    getStudentDetails: async (userId: number): Promise<StudentDetailsResponse> => {
        const response = await apiClient.get<StudentDetailsResponse>(`/admin/students/${userId}`);
        return response.data;
    },
};
