import React from 'react'
import { observer } from 'mobx-react'
import { action, makeObservable, observable, toJS } from 'mobx'
import { useParams } from 'react-router';
import post from '../helpers/post';
import { API } from '../const';
import Task from '../templates/task';

class TaskStore {
    id = 0;
    msg = "";
    list = [];
    constructor() {
        makeObservable(this, {
            id: observable,
            msg: observable,
            list: observable,
            setId: action,
        })
    }

    setId = (id) => {
        if (id !== this.id) {
                this.id = id;
                post(`${API}task/get_by_student`,
                (resp) => {
                    console.log(resp.data);
                    if (resp.data.status === 0) {
                        this.list.replace(resp.data.payload);
                        console.log(toJS(this.list));
                    }
                },
                (error) => {
                    
                },
                {
                    id
                }
            )
        }
    }
}

const taskStore = new TaskStore();

const AssignedTasks = observer(() => {
    return <div>
        {
            taskStore.list.map((task, i) => <Task {...task} key={i} />)
        }
    </div>
});

const Student = observer((props) => {
    taskStore.setId(useParams().id)
    return <div>
        <AssignedTasks />
    </div>
});
export default Student;