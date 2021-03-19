import { Form, postReq } from '../user/form';
import { store } from '../user/store';
import { observer } from 'mobx-react';
import { API } from '../const';

const Auth = observer((props) => {    
    const onChange = (e, field) => {
        store.update(field, e.target.value);
    }
    const inputMap = [
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
        postReq(`${API}admin/auth`);
    }

    const submitLabel = "Войти"

    const form = {
        inputMap,
        submit,
        submitLabel
    }

    return <Form {...form} />
});

export default Auth;