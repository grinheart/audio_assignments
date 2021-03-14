
import { makeObservable, observable, action, computed, toJS } from 'mobx';
import post from '../helpers/post';
import { API } from '../const'

class Store {
    chosen = [];
    full = [];
    selected = 0;
    allChosen = false;
    allStudentsAddedMsg = "";
    constructor(props) {
        makeObservable(this, {
            full: observable,
            chosen: observable,
            allChosen: observable,
            selected: observable,
            allStudentsAddedMsg: observable,

            fullList: computed,
            fullListEnabled: computed,

            //setFullList: action,
            setChosen: action,
            setFull: action,
            addChosen: action,
            setSelected: action,
            removeFromChosen: action,
        });
    }

    get fullList() {
        return this.full
                    .filter(el => !this.chosen.includes(el.id))
                    .sort(
                        (stu1, stu2) => 
                            stu1.name.localeCompare(stu2.name)
                        );
    }

    setChosen = (list) => {
        this.chosen.replace(list);
    }

    setFull = (list) => {
        this.full.replace(list);
        this.selected = this.full[0].id;
    }

    addChosen = (e) => {
        this.chosen.push(this.selected);
        if (this.fullListEnabled) this.selected = this.fullList[0].id;
    }

    removeFromChosen = (id) => {
        this.chosen.replace(this.chosen.filter(el => el !== id));
    }

    setSelected = (e) => {
        this.selected = +e.target.value;
        console.log(this.selected);
    }

    get fullListEnabled() {
        return this.fullList.length;
    }
}
export default Store;