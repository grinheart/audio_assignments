import {React} from 'react';
import { observer } from 'mobx-react';
import {API} from "../const"
import { store } from './store';
import {Form, postReq} from './form';

const Reg = observer((props) => {
    store.redirectIfLogged()

    const onChange = (e, field) => {
        store.update(field, e.target.value);
    }
    const inputMap = [
        {
            type: "text",
            label: "Имя",
            name: "name",
            value: store.name,
            onChange: (e) => onChange(e, "name")
        },
        {
            type: "text",
            label: "Почта",
            name: "email",
            value: store.email,
            onChange: (e) => onChange(e, "email")
        },
        {
            type: "password",
            label: "Пароль",
            name: "pwd",
            value: store.pwd,
            onChange: (e) => onChange(e, "pwd")
        },
    ]

    const submit = () => {
        postReq(`${API}reg`);
    }

    const submitLabel = "Зарегистрироваться"

    const form = {
        inputMap,
        submit,
        submitLabel
    }

    return <Form {...form} />
});

export default Reg;