const { createApp, ref } = Vue

class LocalStorageAdapter {
    constructor() {
        if (!LocalStorageAdapter.instance) {
            LocalStorageAdapter.instance = this;
        }

        return LocalStorageAdapter.instance;
    }

    add(key, value) {
        localStorage.setItem(key, JSON.stringify(value));
    }

    get(key) {
        const item = localStorage.getItem(key);
        return item ? JSON.parse(item) : null;
    }

    remove(key) {
        localStorage.removeItem(key);
    }
}

const ls = new LocalStorageAdapter()

function debounce(func, timeout = 300){
let timer;
    return (...args) => {
        clearTimeout(timer);
        timer = setTimeout(() => { func.apply(this, args); }, timeout);
    };
}
  
const ProductsList = {
    template:
    `
    <div  class="list-container">
        <h2>Product List</h2>
        <a class="checkout-button" href="/checkout">
        <button>Checkout</button>
      </a>
        <input type="text" v-model="search_filter">

        <p>Input value: {{ search_filter }}</p>
            <ul v-if="show_list">
                <li v-for="product in products" :key="product.title">
                    <span>{{ product.title }}</span> - <span>{{ product.price }}</span> 
                    <button class="item-add-button" @click="handleItemClick(product)">Add</button>
                </li>
            </ul>
            <p>{{ total }}</p>
        </div>

    `,
    data() {
        return {
            search_filter: "",
            search_products: true,
            show_list: true,
            total: 0,
            products: [],
            cart_items: []
        }
    },
    mounted() {
        const params = new URLSearchParams(window.location.search);
        const _f = params.get('filter');

        if (_f) {
            this.search_filter = _f;
        }
        fetch(`/api/list?filter=${this.search_filter}`)
            .then(response => response.json())
            .then(data => {
                // Set products data
                this.products = data.body.products;
                this.total = data.body.total
            })
            .catch(error => {
                console.error('Error fetching products:', error);
            });
    },
    watch: {
        search_filter(newVal) {
            this.handleInput(newVal)
        },
    },
    methods: {
        handleInput() {
            // TODO: search every 500 milliseconds
            setTimeout(() => {
                this.search_products = true
            }, 500)
            if(this.search_products == true){
                fetch(`/api/list?filter=${this.search_filter}`)
                    .then(response => response.json())
                    .then(data => {
                        this.products = data.body.products;
                        this.total = data.body.total
                    })
                    .catch(error => {
                        console.error('Error fetching products:', error);
                    });
            }
            this.search_products = false
        },
        handleItemClick(item) {
            this.cart_items = ls.get("cart") == null ? [] : ls.get("cart") 
            this.cart_items.push(item)
            ls.add("cart", this.cart_items)
        }
    },
}

const ProductsListInCart = {
    template:
    `
        <div class="modal" v-show="showSubmitModal">
            <div class="modal-content">
                <span class="close" @click="closeSubmitModal">&times;</span>
                <h2>Confirmation</h2>
                <p>Are you sure you want to submit the form?</p>
                <button class="modal-button" style="background-color: #718472" @click="submitForm">Yes</button>
                <button class="modal-button" style="background-color: #25d428" @click="closeSubmitModal">No</button>
            </div>
        </div>

        <div class="modal" v-show="showErrorModal">
            <div class="modal-content">
                <span class="close" @click="closeErrorModal">&times;</span>
                <h2>Error occured?</h2>
                <button class="modal-button" style="background-color: #718472" @click="closeErrorModal">Ok</button>
            </div>
        </div>

        <div class="form-container">
            <h2>Add New Address</h2>
            <form id="addressForm" @submit.prevent="openSubmitModal">
                <input type="text" v-model="address_form.city" placeholder="City" required>
                <input type="text" v-model="address_form.area"  placeholder="Area" required>
                <input type="number" v-model="address_form.code"  placeholder="Code" required>
                <input type="text" v-model="address_form.street"  placeholder="Street" required>
                <input type="number" v-model="address_form.streetNumber"  placeholder="Street Number" required>
                <input type="hidden" name="csrf" :value="csrf">
                <button type="submit">Submit</button>
            </form>
        </div>
        <div class="list-container">
            <h2>Cart</h2>
            <ul>
                <li v-for="item in cart_items" :key="item.title">
                    <span>{{ item.title }}</span> - <span>{{ item.price }}</span> 
                </li>
            </ul>
        </div>
    `,
    data() {
        return {
            showSubmitModal: false,
            showErrorModal: false,
            csrf: "",
            cart_items: [],
            address_form: {
                city: "",
                area: "",
                code: 0,
                street: "",
                streetNumber: 0
            }
        }
    },
    mounted() {
        this.cart_items = ls.get("cart")
        this.csrf = window._csrf
    },
    methods: {
        openSubmitModal() {
            this.showSubmitModal = true;
          },
        closeSubmitModal() {
            this.showSubmitModal = false;
        },
        openErrorModal() {
            this.showErrorModal = true;
          },
        closeErrorModal() {
            this.showErrorModal = false;
        },
        submitForm() {
            this.closeSubmitModal();
            _products = ls.get("cart")
            _payload = {
                products: _products,
                address: this.address_form,
                csrf: this.csrf
            }
            fetch('/api/order', {
                method: 'POST',
                headers: {
                'Content-Type': 'application/json'
                },
                body: JSON.stringify(_payload)
            })
            .then( response => {
                if (response.ok) {
                    ls.remove("cart")
                    window.location.href = '/thank-you';
                } else {
                    this.openErrorModal()
                }
            })
            .catch( _ => {
                this.openErrorModal()
            });
        }
    }
}

createApp(ProductsList).mount('#product-list')
createApp(ProductsListInCart).mount('#cart-list')