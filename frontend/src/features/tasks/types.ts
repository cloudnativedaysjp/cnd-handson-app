export type Task = {
    id: string;
    title: string;
    description?: string;
    status: string;
    startTimeAt?: string;
    endTimeAt?: string;
    createdAt: string;
    updatedAt: string;
    columnId: string;
    assigneeId: string;
};