import React from 'react'
import { useParams } from 'react-router'
import { action, computed, makeObservable, observable } from 'mobx';
import { observer } from 'mobx-react';
import post from '../helpers/post';
import { API } from '../const';
import Recorder from './recorder'
import axios from 'axios';


const recorder = new Recorder()

class Store {
    id = 0;
    msg = "";
    task = {};
    constructor() {
        makeObservable(this, {
            id: observable,
            msg: observable,
            task: observable,
            setId: action,
            startRecording: action,
        })
    }

    setId = (id) => {
        if (id !== this.id) {
            this.id = id;
                post(`${API}task/check_if_assigned`,
                    (resp) => {
                        if (resp.data.status === 1) {
                        }
                        else {
                                post(`${API}task/get_one`,
                                (resp) => {
                                    console.log(resp.data);
                                    if (resp.data.status === 0) {
                                        Object.assign(this.task, resp.data.payload[0]);   
                                        recorder.setup();
                                    }
                                    else {
                                        this.msg = resp.message
                                    }
                                },
                                (error) => {
                                    this.msg = "Неизвестная ошибка"
                                },
                                {
                                    id
                                }
                            )
                        }
                    },
                    (error) => {
                        this.msg = "Неизвестная ошибка"
                    },
                    {
                        task_id: id,
                    }
                )
        }
    }

    startRecording = () => {
        if (recorder.isRecorderSet) {
            recorder.start();
        }
        else {
            this.msg = "Неизвестная ошибка"
        }
    }

    sendAudio = () => {
        const formData = new FormData();
        formData.append('audio', recorder.audioBlob);
        formData.append('id', this.id);
        const options = {
            withCredentials: true,
            headers:
                {
                    'Content-Type': 'multipart/form-data',
                },
          };
          axios.post(`${API}upload`, formData, options);
    }
}

const store = new Store()

const Task = observer(() => {
    store.setId(useParams().id)
    return <div>
        <p>{store.task.title}</p>
        <p>{store.task.body}</p>
        <button disabled={!recorder.isRecorderSet} onClick={store.startRecording}>Начать</button>
        <button disabled={(!recorder.started || recorder.audioReady)} onClick={recorder.stop}>Закончить</button>
        <button disabled={!recorder.audioReady} onClick={recorder.play}>Прослушать</button>
        <button disabled={!recorder.audioReady} onClick={store.sendAudio}>Отправить</button>
    </div>
});

export default Task;