import { describe, it, expect } from "vitest";
import App from "../App.vue";
import { mount } from "@vue/test-utils";

describe("App", () => {
  it("should render correctly", () => {
    const wrapper = mount(App);
    expect(wrapper.html()).toMatchInlineSnapshot(`
      "<header data-v-f13b4d11=\\"\\">
        <!-- Searcher -->
        <nav class=\\"fixed py-4 bg-white flex flex-row flex-nowrap justify-center items-center p-2 px-4 min-w-full top-nav-shadow\\" data-v-f13b4d11=\\"\\">
          <form class=\\"flex flex-row justify-center w-full\\">
            <ul class=\\"flex flex-row flex-nowrap gap-1 w-10/12\\">
              <li class=\\"w-full\\"><input class=\\"border-gray-900 border-solid border rounded-md p-1 w-full\\" type=\\"store.search\\" name=\\"search_field\\" placeholder=\\"Term for search\\"></li>
              <li><button class=\\"bg-sky-500 hover:bg-sky-700 active:bg-sky-700 focus:outline-none focus:ring focus:ring-sky-300 px-5 py-2 text-sm leading-5 rounded-full font-semibold text-white\\" type=\\"button\\">search</button></li>
            </ul>
          </form>
        </nav>
        <!--==========-->
      </header>
      <div class=\\"p-8 pt-20\\" data-v-f13b4d11=\\"\\"><pre data-v-f13b4d11=\\"\\">
          </pre>
      </div>
      <div class=\\"bg-red-700 fixed top-1/2 top_button\\" data-v-f13b4d11=\\"\\"><button type=\\"button\\" class=\\"btn_pnl_shower rounded-md bg-sky-500 hover:bg-sky-700 active:bg-sky-700 focus:outline-none focus:ring focus:ring-sky-300 px-5 py-2 text-sm leading-5 font-semibold text-white absolute right-0 top-1/2\\" data-v-f13b4d11=\\"\\"><svg xmlns=\\"http://www.w3.org/2000/svg\\" fill=\\"none\\" viewBox=\\"0 0 24 24\\" stroke-width=\\"1.5\\" stroke=\\"currentColor\\" aria-hidden=\\"true\\" class=\\"h-6 w-6\\" data-v-f13b4d11=\\"\\">
            <path stroke-linecap=\\"round\\" stroke-linejoin=\\"round\\" d=\\"M18.75 19.5l-7.5-7.5 7.5-7.5m-6 15L5.25 12l7.5-7.5\\"></path>
          </svg></button></div>"
    `)
  });
});
