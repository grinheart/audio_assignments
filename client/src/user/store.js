import { makeObservable, observable, computed, action } from 'mobx';
import axios from 'axios';
import {API} from '../const';
class Store {
    name = "";
    email = "";
    pwd = "";
    message = "";
    constructor() {
        makeObservable(this, {
          name: observable,
          email: observable,
          pwd: observable,
          raw: computed,
          message: observable,
          update: action,
          setMessage: action,
        });
    }
    get raw() {
        let res = {};
        for (const [key, value] of Object.entries(this)) {
            res[key] = value;
        }
        return res;
    }

    redirectIfLogged = () => {
        const options = {
            withCredentials: true,
            headers:
                {
                    'Content-Type': 'application/x-www-form-urlencoded'
                },
          };
          axios.post(`${API}redirect`, "", options).then(resp => {
              console.log(resp.data);
              if (resp.data.status === 0) {
                  window.location.href = "/"
              }
          })
    }

    update = (name, val) => {
        this[name] = val;
    }

    setMessage = (val) => {
        this.message = val;
    }
}

export const store = new Store()
