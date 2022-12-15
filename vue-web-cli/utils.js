
export const getAuth = () => {
  return {
    headers: {
      Authorization: `Basic ${window.btoa(unescape(encodeURIComponent('admin:Complexpass#123')))}`
    }
  }
}