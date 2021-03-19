import post from '../helpers/post';
import Task from '../templates/task';
import { action, makeObservable, observable } from 'mobx';
import { observer } from 'mobx-react';
import styled from 'styled-components'

class Store {
    tasks = [];
    tasksSet = false;
    msg = "";
    constructor() {
        makeObservable(this, {
            tasks: observable,
            tasksSet: observable,
            msg: observable,
            fetchTasks: action,
        });
    }

    fetchTasks(url) {
        if (this.tasksSet) return;
        post(
            url,
            (resp) => {
                if(resp.data.status === 0) {
                    this.tasks.replace(resp.data.payload);
                }
                else {
                    this.msg = resp.data.message;
                }
            },
            (error => this.msg = "Неизвестная ошибка")
        )
        this.tasksSet = true;
    }
}

const store = new Store();

const Tasks = observer(({url, children, className}) => {
    store.fetchTasks(url);
    return <div className={className}>
        {
            store.tasks.map((task, i) => <Task {...task} key={i}>{children}</Task>)
        }
        {store.msg}
    </div>
});

const StyledTasks = styled(Tasks)`
    display: flex;
    flex-flow: row wrap;
`

export default StyledTasks;