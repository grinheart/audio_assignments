import React from 'react'
import { observer } from 'mobx-react';
import { toJS } from 'mobx';

const AddList = observer((props) => {
    return<select value={props.store.selected} onChange={props.store.setSelected}>
        {
            props.store.fullList
                .map(
                    el =>
                    <option value={el.id} key={el.id}>
                        {el.name} ({el.id})
                    </option>
                    )
        }
    </select>
});

const ChosenList = observer((props) => {
    return <div>
        {
                props.store.chosen.map(id => 
                    <p key={id}>
                        {props.store.full.find(el => el.id === id).name}
                        <button onClick={() => props.store.removeFromChosen(id)}>X</button>
                    </p>)
        }
    </div>
});

const AddingBlock = observer((props) => {
    return <div>
        {props.store.fullListEnabled ?
        <div>
            <p><AddList store={props.store} /></p>
            <button onClick={props.store.addChosen}>Добавить</button>
        </div>
        : props.store.allStudentsAddedMsg}
    </div>
});

const ChoiceBlock = observer((props) => {
    return <div>
        <span>{props.store.title}</span>
        <ChosenList store={props.store} />
    </div>
});

const Choose = observer((props) => {
    return <div>
              <AddingBlock store={props.store} />
              <ChoiceBlock store={props.store} />
          </div>;
});

export default Choose;