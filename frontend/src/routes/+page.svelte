<script>
    const API_ENDPOINT = 'localhost:2000' //will use dns or smth
    let rate = $state(0);
    let text = $state("");
    let isLoading = false;
    let matchWords = $derived(text.match(/\S+/g));
    let wordCount = $derived(matchWords ? matchWords.length : 0)

    async function submit() {
        isLoading = true;

        const data = {
            text: text,
            rate: rate
        };
        
        const response = await fetch(API_ENDPOINT, {
            method: 'POST',
            body: JSON.stringify({data}),
            headers: {
                'Content-Type': 'application/json',
            }
        });
    }

</script>

<p class="logo">SpellWrong</p>

<form class = "inputform">
    <textarea bind:value={text} class = "usertext" placeholder="Paste your text here."></textarea>
    <label>Text length: {wordCount}</label>
    <label>
        <input class="numbers" type="number" bind:value={rate} min="1" max ={wordCount} />
        <!--Super arbitrary, but it makes the thing look nicer-->
        <input class="slider" type="range" bind:value={rate} min="1" max ={Math.round(wordCount/2)} />
    </label>
    <p>Every 1 in {Math.round(wordCount/rate)} words will be misspelled</p>
    <button onclick={submit} class="sendit">Submit</button>
    
</form>

<style>
    .logo {
    font-family: "BBH Sans Bartle", sans-serif;
    font-weight: 400;
    font-style: normal;
    font-size: 40px;
    color: rgb(197, 197, 197);
    display: flex;
    justify-content: center;
    }

    .inputform{
        display: flex;
        align-items: center;
        flex-direction: column;
    }
    
    .usertext{
        min-height: 50vh;
        width: 50%;
        overflow: auto;
        resize: none;
        font-family: Arial, Helvetica, sans-serif;
        border-radius: 10px;
        background-color: rgb(52, 57, 68);
        color: white;
        padding: 10px;
    }

    .sendit{
        font-size: 20px;
        padding: 5px;
        border-radius: 10px;
        color: white;
        background-color: #f26522;
        font-family: "Poppins", sans-serif;
    }
    .slider{
        width: 50vw;
        accent-color: #f26522;
    }
    .numbers{
        color: white;
        border-radius: 5px;
        background-color: rgb(52, 57, 68);
        border: 0px;
    }
</style>