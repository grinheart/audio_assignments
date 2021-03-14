import React from 'react'
import { ChooseStudent, store } from './choosestudent';
import { API } from '../const';
import post from '../helpers/post';
import { action, makeObservable, observable, toJS } from 'mobx';
import { observer } from 'mobx-react'

class Form {
    title = "";
    body = "";
    msg = "";
    constructor() {
        makeObservable(this, {
            title: observable,
            body: observable,
            msg: observable,
            setTitle: action,
            setBody: action,
            setMsg: action,  
          }   
        )   
    }

    setTitle = (e) => {
        this.title = e.target.value
    }

    setBody = (e) => {
        this.body = e.target.value
    }

    setMsg = (msg) => {
        this.msg = msg
    }
}

const f = new Form();

const create = () => {
    post(`${API}task/create`,
    (resp) => {
        if (resp.Status !== 0) {
            f.setMsg(resp.message)
        }
        else {
            f.setMsg('Задание успешно создано')
        }
    },
    (error) => {
        f.setMsg("Неизвестная ошибка")
    },
    {
        title: f.title,
        body: f.body,
        assign_to: toJS(store.chosen).join(","),
    })
}

const CreateTask = observer(() => {
    return <div>
        <input type="text" name="title" value={f.title} onChange={f.setTitle} />
        <textarea name="body" value={f.body} onChange={f.setBody} />
        <ChooseStudent />
        <button onClick={create}>Создать</button>
        <span>{f.msg}</span>
    </div>
});

export default CreateTask;