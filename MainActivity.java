package com.example.yasmi.unirest;

import android.content.Intent;
import android.os.Bundle;
import android.support.v7.app.AppCompatActivity;
import android.util.Log;
import android.view.View;
import android.widget.Button;

import org.json.JSONException;
import org.json.JSONObject;

import java.io.IOException;

import okhttp3.Call;
import okhttp3.Callback;
import okhttp3.OkHttpClient;
import okhttp3.Request;
import okhttp3.Response;

public class MainActivity extends AppCompatActivity {
    private  static final String TAG = "MainActivity";

    private  String auth = "unknown";
    private String welcome="";


    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);

        Button welcomeBtn= (Button) findViewById(R.id.welcomeBtn);
        welcomeBtn.setOnClickListener(new View.OnClickListener() {
            Intent ChatActivity=new Intent(getApplicationContext(),Chat.class);
            public void onClick(View v) {

                OkHttpClient client = new OkHttpClient();

                Request request = new Request.Builder()
                        .url("https://safe-everglades-93552.herokuapp.com/welcome")
                        .get()
                        .build();


                client.newCall(request).enqueue(new Callback() {
                    @Override
                    public void onFailure(Call call, IOException e) {

                        Log.i(TAG,e.getMessage()) ;
                    }

                    @Override
                    public void onResponse(Call call, Response response) throws IOException {

                        String jsonStr=response.body().string(); //.substring(0+12,0+11+44);
                        Log.i(TAG, "onResponse: " + jsonStr);

                       try{
                           JSONObject jsonObj = new JSONObject(jsonStr);
                           welcome = jsonObj.getString("message");
                          Log.i(TAG,welcome);
                            auth = jsonObj.getString("uuid");
                           Log.d(TAG,auth);
                       }catch (final JSONException e) {
                           Log.e(TAG, "Json parsing error: " + e.getMessage());

                       }

                        ChatActivity.putExtra("com.example.yasmi.unirest.auth",auth);
                        ChatActivity.putExtra("com.example.yasmi.unirest.welcome",welcome);
                        startActivity(ChatActivity);
                    }
                });



            }
        });






            }

}