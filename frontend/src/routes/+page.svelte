<script>
    let mistakes = $state(0);
    let text = $state("");
    let isLoading = $state(false);
    let resultRecieved = $state(false);
    let error = null;
    let responseMessage = $state('');
    let matchWords = $derived(text.match(/\S+/g));
    let wordCount = $derived(matchWords ? matchWords.length : 0);

    async function submit() {
        isLoading = true; // State is set to true when the request starts
        error = null;     
        responseMessage = '';

        const data = {
            text: text,
            mistakes: mistakes,
            wordCount: wordCount
        };

        const bodyContent = JSON.stringify(data); 

        try {
            const response = await fetch("/submit", {
                method: 'POST',
                body: bodyContent,
                headers: {
                    'Content-Type': 'application/json',
                }
            });

            // 1. Check if the HTTP status code indicates success (200-299)
            if (response.ok) {
                // 2. Parse the successful JSON response from the server
                const result = await response.json();
                resultRecieved = true;
                responseMessage = result.responseData.receivedText;
                
            } else {
                // 3. Handle specific HTTP error status codes (e.g., 400, 500)
                const errorData = await response.text(); // Read the error body
                throw new Error(`HTTP Error ${response.status}: ${errorData}`);
            }
        } catch (e) {
            // 4. Handle network errors (e.g., server offline, connection failed)
            console.error('Submission failed:', e);
            error = e.message;
        } finally {
            // 5. This block runs after try or catch, ensuring isLoading is reset
            
            isLoading = false;
            
        }
    }

</script>

<p class="logo">SPELLWRONG</p>

{#if !isLoading && !resultRecieved}
    <form class = "inputform">
        <textarea bind:value={text} class = "usertext" placeholder="Paste your text here."></textarea>
        <label>Word count: {wordCount}</label>
        <label>
            <input class="numbers" type="number" bind:value={mistakes} min="1" max ={wordCount} />
            <input class="slider" type="range" bind:value={mistakes} min="1" max ={wordCount} />
        </label>
        {#if mistakes == wordCount}
            <p>We will try to misspell every word</p>
        {:else}
            <p>{mistakes} {mistakes == 1 ? "word" : "words"} will be misspelled, which is equivalent to every 1 in {Math.round(wordCount/mistakes)} words. </p>
        {/if}
        <button onclick={submit} class="sendit">Submit</button>
        
    </form>
{:else if isLoading}
    <p class="loading">Loading...</p>
{:else if !isLoading && resultRecieved}
    <div class="presentResult">
        <p>Below is your new text</p>
        <p class = "output">{responseMessage}</p>
    </div>
    
{/if}



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

    .output{
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

    .presentResult{
        display: flex;
        flex-direction: column;
        align-items: center;
    }
    
    .loading{
        font-size: x-large;
        position: absolute;
        margin-left: auto;
        margin-right: auto;
        left: 0;
        right: 0;
        text-align: center;    
    }
</style>