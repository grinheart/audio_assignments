import {React} from 'react';
import post from '../helpers/post';
import { API } from '../const';
const Logout = () => {
    post(`${API}logout`);
    window.location.href = "/auth";
    return <div></div>
}
export default Logout