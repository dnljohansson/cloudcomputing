//server that sends the request to dispatch internally

import {json} from '@sveltejs/kit'
import {API_ENDPOINT} from '$env/static/private'

export async function POST({request}) {
    const requestData = await request.json();

    const response = await fetch(API_ENDPOINT, {
        method: 'POST',
        body: JSON.stringify(requestData),
        headers: {
            'Content-Type': 'application/json'
        }
    });

    const responseData = await response.json()

    return json({responseData}, {status: response.status});

}