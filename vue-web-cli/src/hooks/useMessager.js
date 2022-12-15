import { ref } from 'vue'

export const useMessager = () => {
  const message = ref("use Message")

  return {
    message
  }
}