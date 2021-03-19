import React from 'react'
import { observer } from 'mobx-react'
import { API } from '../const'
import ReactAudioPlayer from 'react-audio-player'
import styled from 'styled-components'

const Task = observer(({id, title, body, audio, children, className}) => {
    if (!children) {
        children = () => {};
    }

    const Title = styled.p`
        font-size: 30px;
    `;

    return <div className={className}>
        <Title>{title}</Title>
        <p>{body}</p>
        <div>
            {
                audio.map((path, i) => 
                    <div><ReactAudioPlayer
                        src={API + path}
                        crossOrigin="anonymous"
                        controls
                        key={i}
                        />
                    </div>
                )
            }
        </div>
        <div>{children({id, title, body})}</div>
    </div>
})

const StyledTask = styled(Task)`
    background-color: #98FF98;
    width: calc((100vw - 70px) / 3);
    margin-left: 10px;
    margin-top: 10px;
    padding: 5px;
    border-radius: 10px;
`;

export default StyledTask;