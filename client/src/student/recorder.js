import { action, computed, makeObservable, observable } from 'mobx';

class Recorder {
    mediaRecorder = 0;
    audioChunks = [];
    audioBlob = 0;
    audio = 0;
    started = false;
    audioReady = false;
    constructor() {
        makeObservable(this, {
            start: action,
            isRecorderSet: computed,
            started: observable,
            mediaRecorder: observable,
            audio: observable,
            audioReady: observable,
            setAudioReady: action,
            setup: action,
        });
    }

    setup = () => {
        navigator.mediaDevices.getUserMedia({ audio: true })
            .then(stream => {
            this.mediaRecorder = new MediaRecorder(stream);
            this.mediaRecorder.addEventListener("stop", () => {
                this.audioBlob = new Blob(this.audioChunks);
                const audioUrl = URL.createObjectURL(this.audioBlob);
                this.audio = new Audio(audioUrl);
                this.setAudioReady();
            });
            },
            error => {
                console.log(error)
            }
        );
    }

    setAudioReady() {
        this.audioReady = true;
    }

    start = () => {
        if (this.mediaRecorder) {
            this.audioChunks = [];
            this.mediaRecorder.start();
            this.mediaRecorder.addEventListener("dataavailable", event => {
            this.audioChunks.push(event.data);
            });
            this.started = true;
            this.audioReady = false;
            return true;
        }
        else {
            return false;
        }
    }

    stop = () => {
        this.mediaRecorder.stop();
    }

    get isRecorderSet() {
        return this.mediaRecorder;
    }

    play = () => {
        this.audio.play()
    }
}

export default Recorder;