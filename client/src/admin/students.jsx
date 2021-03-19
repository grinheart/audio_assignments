import React from 'react'
import { observer } from 'mobx-react'
import { action, makeObservable, observable, toJS } from 'mobx'
import post from '../helpers/post';
import { API } from '../const';
import axios from 'axios';

class StudentsStore {
    id = 0;
    list = [];
    listSet = false;
    msg = "";
    constructor() {
        makeObservable(this, {
            id: observable,
            list: observable,
            listSet: observable,
            setList: action,
            msg: observable,
        });
    }

    setList() {
        if (this.listSet) return;
        axios.get(`${API}admin_redirect`).then((resp) => {
        })
        post(
            `${API}students`,
            (resp) => {
                console.log(resp);
                if (resp.data.status === 0) {
                    console.log("at 0");
                    this.listSet = true;
                    this.list.replace(resp.data.payload);
                }
            },
            (error) => {
                this.msg = "Неизвестная ошибка";
            }

        )
    }
}

const store = new StudentsStore();

const Students = observer(() => {
    store.setList();
    return <div>
        <div>
        {
            store.list.map(stu => {
                return <p><a href={`/admin/student/${stu.id}`}>{stu.name}</a></p>
            })
        }
        </div>
        <p>{store.msg}</p>
    </div>
});

export default Students;