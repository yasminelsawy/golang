package com.example.yasmi.unirest;

import android.content.Intent;
import android.net.Uri;
import android.os.Bundle;
import android.support.v7.app.AppCompatActivity;
import android.util.Log;
import android.view.View;
import android.widget.Button;
import android.widget.EditText;
import android.widget.ListView;
import android.widget.TextView;

import com.github.library.bubbleview.BubbleTextView;

import org.json.JSONException;
import org.json.JSONObject;

import java.io.IOException;

import okhttp3.Call;
import okhttp3.Callback;
import okhttp3.MediaType;
import okhttp3.OkHttpClient;
import okhttp3.Request;
import okhttp3.RequestBody;
import okhttp3.Response;

public class Chat extends AppCompatActivity {

    public static  ChatMessageHistory history =new ChatMessageHistory();
    boolean wlcmHistory=false;
    boolean resetHist=false;
    String authrset="";

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        final Intent welcome=getIntent();
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_chat);
        TextView authTextView=(TextView) findViewById(R.id.auth);
        String auth=welcome.getStringExtra("com.example.yasmi.unirest.auth");
        authTextView.setText(auth);
        TextView welcomeMessage=(BubbleTextView) findViewById(R.id.wlcmMessage);
        String welcomee=welcome.getStringExtra("com.example.yasmi.unirest.welcome");
        welcomeMessage.setText(welcomee);
        if(!wlcmHistory){
            history.setMessageServerHistory(welcomee);
            wlcmHistory=true;
            Log.e("BOOLEAN",""+ wlcmHistory);
        }




        Button sendBtn=(Button)findViewById(R.id.sendBtn);
        sendBtn.setOnClickListener( new View.OnClickListener(){

            public void onClick(View v){
                TextView welcomeMessage=(BubbleTextView) findViewById(R.id.wlcmMessage);
                welcomeMessage.setVisibility(View.INVISIBLE);
                EditText userInput=(EditText)findViewById(R.id.userInput);
                Button newButton=(Button) findViewById(R.id.newButton);
                newButton.setVisibility(View.INVISIBLE);
                Button usedbutton=(Button)findViewById(R.id.usedbutton);
                usedbutton.setVisibility(View.INVISIBLE);
                String input=userInput.getText().toString();

                userInput.setText("");

                if(input!="Send message ..."){
                    Intent welcome=getIntent();
                    history.setMessageUserHistory(input);
                    Log.e("Size Of server before",""+history.getSizeOfServerHistory());
                    Log.e("Size Of User",""+history.getSizeOfUserHistory());

                    Log.e("Size Of server after",""+history.getSizeOfServerHistory());
                    OkHttpClient client = new OkHttpClient();
                    JSONObject main = new JSONObject();
                    try {
                        main.put("message", input);
                        Log.e("JSON",main.toString());
                    }catch (JSONException e){
                        Log.e("FAILLL", e.getMessage());
                    }

                    String auth="";
                    if(resetHist){
                        auth =authrset;
                    }else{

                        auth=welcome.getStringExtra("com.example.yasmi.unirest.auth");
                    }
                    MediaType JSON = MediaType.parse("application/json; charset=utf-8");
                    RequestBody body = RequestBody.create(JSON, main.toString());
                    Request request = new Request.Builder()
                            .url("https://safe-everglades-93552.herokuapp.com/chat")
                            .post(body)
                            .addHeader("Authorization", auth)
                            .build();
                    client.newCall(request).enqueue(new Callback() {
                        public void onFailure(Call call, IOException e) {
                            Log.e("Response", e.getMessage());
                        }
                        @Override
                        public void onResponse(Call call, Response response) throws IOException {

                            String jsonStr=response.body().string();

                            try{
                                JSONObject jsonObj = new JSONObject(jsonStr);

                                final String data = jsonObj.getString("message");
                                history.setMessageServerHistory(data);
                                Log.i("RESPONSE",data);
                                final ChatAdapter newHistory =new ChatAdapter(getApplicationContext(),history.getServerMessage(),history.getUserMessage());


                                if (data.contains("https")==true){

                                    final  String data1=data.substring(44,data.length());
                                    Log.i("ANAA HENAAAAAA", data1);
                                    Chat.this.runOnUiThread(new Runnable() {
                                        @Override
                                        public void run() {
                                            final ListView chatHistory =(ListView) findViewById(R.id.historyList);
                                            chatHistory.setAdapter(newHistory);


                                            Intent ii = new Intent(Intent.ACTION_VIEW);
                                            ii.setData(Uri.parse(data1));
                                            ii.setFlags(Intent.FLAG_ACTIVITY_NEW_TASK | Intent.FLAG_ACTIVITY_CLEAR_TASK);
                                            startActivity(ii);

                                        }
                                    });


                                }else{
                                    runOnUiThread(new Runnable() {
                                        @Override
                                        public  void run() {
                                            Log.e("Size Of server",""+history.getSizeOfServerHistory());
                                            Log.e("Size Of User",""+history.getSizeOfUserHistory());
                                            ListView chatHistory =(ListView) findViewById(R.id.historyList);
                                            chatHistory.setAdapter(newHistory);

                                        }
                                    });
                                }




                            }catch (final JSONException e) {
                                Log.i("ERROR",e.getMessage());

                            }

                        }
                    });

                }
            }





        });


        //the quick action button to send new 3la tool


        final Button newButton=(Button)findViewById(R.id.newButton);
        newButton.setOnClickListener( new View.OnClickListener(){

            public void onClick(View v){
                TextView welcomeMessage=(BubbleTextView) findViewById(R.id.wlcmMessage);
                welcomeMessage.setVisibility(View.INVISIBLE);
                newButton.setVisibility(View.INVISIBLE);
                Button usedButton =(Button) findViewById(R.id.usedbutton);
                usedButton.setVisibility(View.INVISIBLE);

                EditText userInput=(EditText)findViewById(R.id.userInput);

                String input=userInput.getText().toString();

                userInput.setText("");

                if(input!="Send message ..."){
                    Intent welcome=getIntent();
                    history.setMessageUserHistory("new");

                    OkHttpClient client = new OkHttpClient();
                    JSONObject main = new JSONObject();
                    try {
                        main.put("message", "new");
                        Log.e("JSON",main.toString());
                    }catch (JSONException e){
                        Log.e("FAILLL", e.getMessage());
                    }

                    String auth="";
                    if(resetHist){
                        auth =authrset;
                    }else{

                        auth=welcome.getStringExtra("com.example.yasmi.unirest.auth");
                    }

                    MediaType JSON = MediaType.parse("application/json; charset=utf-8");
                    RequestBody body = RequestBody.create(JSON, main.toString());
                    Request request = new Request.Builder()
                            .url("https://safe-everglades-93552.herokuapp.com/chat")
                            .post(body)
                            .addHeader("Authorization", auth)
                            .build();
                    client.newCall(request).enqueue(new Callback() {
                        public void onFailure(Call call, IOException e) {
                            Log.e("Response", e.getMessage());
                        }
                        @Override
                        public void onResponse(Call call, Response response) throws IOException {

                            String jsonStr=response.body().string();

                            try{
                                JSONObject jsonObj = new JSONObject(jsonStr);

                                final String data = jsonObj.getString("message");
                                history.setMessageServerHistory(data);
                                if(history.getServerMessage().size()==0){
                                    Log.e("YALHWIII","EL ARRAY FADI BTA# EL SERVER");
                                }
                                if(history.getUserMessage().size()==0){
                                    Log.e("YALHWIII","EL ARRAY FADI BTA# EL USER");
                                }
                                Log.i("RESPONSE",data);

                                final ChatAdapter newHistory =new ChatAdapter(getApplicationContext(),history.getServerMessage(),history.getUserMessage());
                                runOnUiThread(new Runnable() {
                                    @Override
                                    public  void run() {

                                        newButton.setVisibility(View.INVISIBLE);
                                        ListView chatHistory =(ListView) findViewById(R.id.historyList);
                                        chatHistory.setAdapter(newHistory);


                                    }
                                });


                            }catch (final JSONException e) {
                                Log.i("ERROR",e.getMessage());

                            }

                        }
                    });

                }
            }





        });



        final Button usedbutton=(Button)findViewById(R.id.usedbutton);
        usedbutton.setOnClickListener( new View.OnClickListener(){

            public void onClick(View v){
                TextView welcomeMessage=(BubbleTextView) findViewById(R.id.wlcmMessage);
                welcomeMessage.setVisibility(View.INVISIBLE);
                EditText userInput=(EditText)findViewById(R.id.userInput);
                String input=userInput.getText().toString();
                userInput.setText("");
                usedbutton.setVisibility(View.INVISIBLE);
                newButton.setVisibility(View.INVISIBLE);



                if(input!="Send message ..."){
                    Intent welcome=getIntent();
                    history.setMessageUserHistory("used");

                    OkHttpClient client = new OkHttpClient();
                    JSONObject main = new JSONObject();
                    try {
                        main.put("message", "used");
                        Log.e("JSON",main.toString());
                    }catch (JSONException e){
                        Log.e("FAILLL", e.getMessage());
                    }
                    String auth="";
                    if(resetHist){
                        auth =authrset;
                    }else{

                        auth=welcome.getStringExtra("com.example.yasmi.unirest.auth");
                    }
                    MediaType JSON = MediaType.parse("application/json; charset=utf-8");
                    RequestBody body = RequestBody.create(JSON, main.toString());
                    Request request = new Request.Builder()
                            .url("https://safe-everglades-93552.herokuapp.com/chat")
                            .post(body)
                            .addHeader("Authorization", auth)
                            .build();
                    client.newCall(request).enqueue(new Callback() {
                        public void onFailure(Call call, IOException e) {
                            Log.e("Response", e.getMessage());
                        }
                        @Override
                        public void onResponse(Call call, Response response) throws IOException {

                            String jsonStr=response.body().string();

                            try{
                                JSONObject jsonObj = new JSONObject(jsonStr);

                                final String data = jsonObj.getString("message");
                                history.setMessageServerHistory(data);
                                if(history.getServerMessage().size()==0){
                                    Log.e("YALHWIII","EL ARRAY FADI BTA# EL SERVER");
                                }
                                if(history.getUserMessage().size()==0){
                                    Log.e("YALHWIII","EL ARRAY FADI BTA# EL USER");
                                }
                                Log.i("RESPONSE",data);

                                final ChatAdapter newHistory =new ChatAdapter(getApplicationContext(),history.getServerMessage(),history.getUserMessage());
                                runOnUiThread(new Runnable() {
                                    @Override
                                    public  void run() {

                                        newButton.setVisibility(View.INVISIBLE);
                                        ListView chatHistory =(ListView) findViewById(R.id.historyList);
                                        chatHistory.setAdapter(newHistory);


                                    }
                                });


                            }catch (final JSONException e) {
                                Log.i("ERROR",e.getMessage());

                            }

                        }
                    });

                }
            }





        });

        final Button reset=(Button) findViewById(R.id.reset);

        reset.setOnClickListener( new View.OnClickListener(){

            public void onClick(View v){
                resetHist=!resetHist;
                TextView welcomeMessage=(BubbleTextView) findViewById(R.id.wlcmMessage);
                // welcomeMessage.setVisibility(View.INVISIBLE);
                EditText userInput=(EditText)findViewById(R.id.userInput);
                // TextView UserMessage =(BubbleTextView) findViewById(R.id.userMessage);
                // String input=userInput.getText().toString();
                // UserMessage.setText(input);
                userInput.setText("");
                usedbutton.setVisibility(View.VISIBLE);
                newButton.setVisibility(View.VISIBLE);

                // ChatAdapter newHistory =new ChatAdapter(getApplicationContext());
                //ListView chatHistory =(ListView) findViewById(R.id.historyList);
                //chatHistory.setAdapter(newHistory);


                //Log.d("Response", "GOWAAA");
                //if(input!="Send message ..."){
                // Intent welcome=getIntent();
                // history.setMessageUserHistory("used");

                OkHttpClient client = new OkHttpClient();
                JSONObject main = new JSONObject();
                //  try {
                //    main.put("message", "used");
                //  Log.e("JSON",main.toString());
                //}catch (JSONException e){
                //  Log.e("FAILLL", e.getMessage());
                //}
                //  Intent welcome=getIntent();
                // String auth=welcome.getStringExtra("com.example.yasmi.unirest.auth");
                MediaType JSON = MediaType.parse("application/json; charset=utf-8");
                //  RequestBody body = RequestBody.create(JSON, main.toString());
                Request request = new Request.Builder()
                        .url("https://safe-everglades-93552.herokuapp.com/welcome")
                        .get()
                        //.addHeader("Authorization", auth)
                        .build();
                client.newCall(request).enqueue(new Callback() {
                    public void onFailure(Call call, IOException e) {
                        Log.e("Response", e.getMessage());
                    }
                    @Override
                    public void onResponse(Call call, Response response) throws IOException {

                        String jsonStr=response.body().string();

                        try{
                            JSONObject jsonObj = new JSONObject(jsonStr);

                            final String data = jsonObj.getString("message");
                            authrset = jsonObj.getString("uuid");
                            resetHist=true;
                            history.reset();
                           // history.setMessageServerHistory(data);
                            if(history.getServerMessage().size()==0){
                                Log.e("YALHWIII","EL ARRAY FADI BTA# EL SERVER");
                            }
                            if(history.getUserMessage().size()==0){
                                Log.e("YALHWIII","EL ARRAY FADI BTA# EL USER");
                            }
                            Log.i("RESPONSE",data);

                            final ChatAdapter newHistory =new ChatAdapter(getApplicationContext(),history.getServerMessage(),history.getUserMessage());
                            runOnUiThread(new Runnable() {
                                @Override
                                public  void run() {
                                    //TextView welcomeMessage=(BubbleTextView) findViewById(R.id.wlcmMessage);
                                    //welcomeMessage.setText(data);
                                    // newButton.setVisibility(View.INVISIBLE);
                                    //ListView chatHistory =(ListView) findViewById(R.id.historyList);
                                    //chatHistory.setAdapter(newHistory);
                                    TextView auth=(TextView) findViewById(R.id.auth);
                                    auth.setText(authrset);


                                }
                            });


                        }catch (final JSONException e) {
                            Log.i("ERROR",e.getMessage());

                        }

                    }
                });

            }
            // }





        });


    }
/* it's the context of current state of the application/object. It lets newly-created objects understand what has been going on.
 Typically you call it to get information regarding another part of your program (activity and package/application).
 */


}
