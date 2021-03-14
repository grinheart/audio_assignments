import React from 'react'
import { observer } from 'mobx-react';
import post from '../helpers/post';
import Store from './store';
import { API } from '../const'
import Choose from './choose'

class St extends Store {
    constructor() {
        super();
        post(`${API}students`, (resp) => {
            this.setFull(resp.data.payload);
            this.selected = this.full[0].id;
        });
        this.allStudentsAddedMsg = "Все студенты добавлены в список"
    }
}

export const store = new St()

export const ChooseStudent = observer(() => {
    return <Choose store={store} />
});
