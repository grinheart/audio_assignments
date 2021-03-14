import { React } from 'react';
import { observer } from 'mobx-react';
import { store } from './store';
import post from '../helpers/post'

const Form = observer((props) => {
    return  <div>
    {
        props.inputMap.map((attrs, i) => 
            <p key={i}> <span>{attrs.label}</span> <input {...attrs} /></p>
        )
    }
    <p><button onClick={props.submit}>{props.submitLabel}</button></p>
    <p>{store.message}</p>
</div>;
});

const postReq = (url) => {
    post(url,
    (resp) => 
    {
        if (resp.data.status === 0) {
          window.location.href = "/";
        }
        else if (resp.data.status === 1) {
            store.setMessage(resp.data.message);
        }
        else {
            store.setMessage("Неизвестная ошибка")
        }
        console.log(resp.data)
    },
    (error) => {
        store.setMessage("Неизвестная ошибка")
    }, store.raw)
}

export {Form, postReq};