'use strict'

let showSpinner = () => {
    document.getElementById('loaderC').style.display = 'block';
    document.getElementById('loader').style.display = 'block';
}

let hideSpinner = () => {
    document.getElementById('loaderC').style.display = 'none';
    document.getElementById('loader').style.display = 'none';
}

let errorToast = (message) => {
    Swal.fire({
        position: 'center',
        icon: 'error',
        text: message,
        showConfirmButton: false,
        timer: 3000
    })
}

let deleteToast = (title, data, callback) => {
    Swal.fire({
        position: 'center',
        icon: 'warning',
        text: `آیا از حذف ${title} مطمئن هستید؟`,
        showCancelButton: true,
        cancelButtonText: 'خیر',
        cancelButtonColor: '#d33',
        confirmButtonColor: '#3085d6',
        confirmButtonText: 'بله',
    }).then((result) => {
        if (result.value) {
            callback(data)
        }
    })
}

let successToastWithReload = (message) => {
    Swal.fire({
        position: 'center',
        icon: 'success',
        text: message,
        showConfirmButton: false,
        timer: 3000
    }).then((result) => {
        location.reload()
    })
}

class Ajax {
    constructor(url, data, method = 'POST') {
        this.method = method
        this.url = url
        this.data = data
        this.xhr = new XMLHttpRequest()
        this.xhr.responseType = 'json'
    }

    post(callback) {
        showSpinner()

        this.xhr.open(this.method, this.url);
        this.xhr.setRequestHeader('Content-type', 'application/json; charset=utf-8');
        this.xhr.send(JSON.stringify(this.data));
        this.xhr.onload = () => {
            hideSpinner()

            if ((this.xhr.status < 500 || this.xhr.status >= 600) &&
                (this.xhr.status < 400 || this.xhr.status > 405)) {
                callback ? callback(this.xhr) : null
            } else {
                let message;
                switch (this.xhr.status) {
                    case 400:
                        message = 'درخواست شما نامعتبر است.';
                        break;
                    case 403:
                        message = 'شما مجوز انجام این کار را ندارید.'
                        break
                    case 404:
                        message = 'آدرس سرور نامعتبر است.';
                        break;
                    default:
                        message = 'خطای سرور، در اسرع وقت با اپراتور تماس بگیرید.';
                        break;
                }
                errorToast(message)
            }
        }
        this.xhr.onerror = () => {
            hideSpinner()
            errorToast('خطایی در ارسال درخواست پیش آمده است.')
        }
    }
}