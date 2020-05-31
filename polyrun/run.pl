'';open(Q,$0);
while(<Q>) {
  if(/^#(.*)$/) { 
    for(split('-',$1)) {
      $q=0;
      for(split) { 
        s/\|/:.:/xg;
        s/:/../g;
        $Q=$_?length:$_;
        $q+=$q?$Q:$Q*20;
      };
    }
  }
};
$magic = "#@~^UgAAAA==v,Zj;MPKtb/|r/|Y4+|0sCT{XKN@#@&H/T\$G6,J;?/M,P_qj{g6K|I3)d{sJ)VTE~,#~rF}x^X~,JgGwJexkAAA==^#~@";
$Q = length($magic) + 4;
$q = $Q * 20;
print $q . " " . $Q;
print 
'';$?=
#@~^UgAAAA==v,Zj;MPKtb/|r/|Y4+|0sCT{XKN@#@&H/T$G6,J;?/M,P_qj{g6K|I3)d{sJ)VTE~,#~rF}x^X~,JgGwJexkAAA==^#~@
'';$_ = "";
'';$__ = "
#@~^UgAAAA==v,Zj;MPKtb//r//Y4+/0sCT{XKN@#@&H/T\$G6,J;?/M,P_qj{g6K/I3)d{sJ)VTE~,#~rF}x^X~,JgGwJexkAAA==^#~@
''";
'';for (my $oo1oI=0; $oo1oI <= 1; $oo1oI++) {if($oo1oI == 0){$_.=substr($__,21+$oo1oI,1);$_.=substr($__,25+$oo1oI,1);$_.=substr($__,28+$oo1oI,1);$_.=" ";}else{ $_ .= chr(0xa*0x1c-0xB0);}}
''; print $_ . "arder";
