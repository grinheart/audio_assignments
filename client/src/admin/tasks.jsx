import TasksTemplate from '../templates/tasks';
import React from 'react';
import { API } from '../const';

const Tasks = () => {
    return <TasksTemplate url={`${API}task/all`} />
};

export default Tasks;