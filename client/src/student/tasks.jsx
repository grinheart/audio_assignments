import TasksTemplate from '../templates/tasks';
import React from 'react';
import { API } from '../const';
import styled from 'styled-components';

const StartButton = ({id}) => {
    return <button onClick={() => window.location.href = `/task/${id}`}>Приступить</button>
}

const Tasks = () => {
    console.log("task")
    return <TasksTemplate url={`${API}task/get_for_student`} children={StartButton}>
    </TasksTemplate>
};

export default Tasks;