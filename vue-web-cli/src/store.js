import { reactive } from 'vue'
import { searchEmail } from './services/EmailService'

export const store = reactive({
  open: false,
  search: '',
  lastSearch: '',
  emails: [],
  emailsTotal: 1,
  currentPage: 1,
  selectedEmailID: '',
  selectedEmailInfo: undefined,
  async onSearch() {
    try {
      if (this.lastSearch != this.search) {
        this.emailsTotal = 1
        this.currentPage = 1
      }
      const { _emails, _total } = await searchEmail(
        this.search,
        {
          from: this.currentPage,
          size: 20,
          count: this.emailsTotal
        }
      )
      this.lastSearch = this.search
  
      //removing unuseful characters
      for (let mail of _emails) {
        for (let d in mail._source) {
          mail._source[d] = mail._source[d].trim()
        }
      }
  
      this.emails = _emails
      this.emailsTotal = _total
  
      this.open = true
    } catch (error) {
      alert(error)
    }
  }
})