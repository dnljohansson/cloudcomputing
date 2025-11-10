//server that sends the request to dispatch internally

import {json} from '@sveltejs/kit'
import {env} from '$env/dynamic/private'

export async function POST({request}) {
    const requestData = await request.json();
    const url = env.API_ENDPOINT;
    const response = await fetch(url, {
        method: 'POST',
        body: JSON.stringify(requestData),
        headers: {
            'Content-Type': 'application/json'
        }
    });

    const responseData = await response.json()

    return json({responseData}, {status: response.status});

}