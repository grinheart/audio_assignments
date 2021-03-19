import React from 'react'
import { observer } from 'mobx-react';
import post from '../helpers/post';
import Store from './store';
import { API } from '../const'
import Choose from './choose'
import { action, makeObservable } from 'mobx';

class St extends Store {
    studentsSet = false;
    constructor() {
        super();
        this.allStudentsAddedMsg = "Все студенты добавлены в список"
        this.title = "Выбранные студенты"
        makeObservable(this, {
            fetchStudents: action,
        })
    }

    fetchStudents() {
        if (this.studentsSet) return;
        post(`${API}students`, (resp) => {
            this.setFull(resp.data.payload);
            this.selected = this.full[0].id;
        });
        this.studentsSet = true;
    }
}

export const store = new St()

export const ChooseStudent = observer(() => {
    store.fetchStudents();
    return <Choose store={store} />
});
