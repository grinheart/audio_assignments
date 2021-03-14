import axios from 'axios'
import qs from 'qs'
const post = (url, atSuccess, atError, body, headers) => {
    const options = {
        withCredentials: true,
        headers:
            {
                'Content-Type': 'application/x-www-form-urlencoded',
                ...Headers
            },
      };
      axios.post(url, qs.stringify(body), options).then(atSuccess, atError)
}

export default post;