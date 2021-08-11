const state = {
    title: 'Top Stories'
};
export const getters = {
    title: (theState) => {
        return theState.title;
    }
};
const mutations = {
    setTitle(theState, title) {
        theState.title = title;
    }
};
export const actions = {
    changeTitle({ commit }, title) {
        commit('setTitle', title);
    }
};
export const topToolbar = {
    namespaced: true,
    state,
    getters,
    mutations,
    actions
};
//# sourceMappingURL=topToolbar.js.map