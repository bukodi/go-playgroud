import Vue from 'vue';
import Vuex from 'vuex';
import { topToolbar } from './modules/topToolbar';
Vue.use(Vuex);
const store = {
    modules: {
        topToolbar
    }
};
export default new Vuex.Store(store);
//# sourceMappingURL=index.js.map