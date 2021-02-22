export default function ({ $axios, redirect }) {
    $axios.onRequest(config => {
        console.log('Making request to ' + config.url)
      })
    $axios.onError(error => {
      if(error.response.status === 500) {
        redirect('/Error500')
      }
    })
  }
  